package apps

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nano-ci-cd/config"
	db "github.com/nano-ci-cd/db"
	types "github.com/nano-ci-cd/types"
	"github.com/nano-ci-cd/util"
)

const (
	configFolder          = ".nano-cicd"
	latestShaShortCommand = "git rev-parse --short HEAD"
	baseContainerFolder   = "/app"
)

type SingleBuildContext struct {
	Context    *context.Context
	Db         *db.AppsDb
	LogsWriter *AppLogsWriter
}

func Build(buildContext context.Context, db *db.AppsDb, logsChan chan string) (SingleBuildContext, error) {
	appWriter := &AppLogsWriter{
		Logs:     "",
		LogsChan: logsChan,
	}

	bContext := SingleBuildContext{
		Context:    &buildContext,
		Db:         db,
		LogsWriter: appWriter,
	}

	os.Chdir(baseContainerFolder)

	println(fmt.Sprintf("Build started at %v", time.Now()))

	err := bContext.cloneRepo(buildContext)

	if err != nil {
		return bContext, err
	}

	err = bContext.prepareEnvAndBuildArguments(buildContext, db)

	if err != nil {
		return bContext, err
	}

	err = bContext.runPreBuildScript()

	if err != nil {
		return bContext, err
	}

	err = bContext.runBuildScript()

	if err != nil {
		return bContext, err
	}

	err = loadBase64EnvFileIntoEnv(buildContext, db)

	if err != nil {
		return bContext, err
	}

	err = bContext.runPostBuildScript()

	if err != nil {
		return bContext, err
	}

	err = bContext.executeDockerComposeFileIfConfigured(buildContext)

	if err != nil {
		return bContext, err
	}

	println(fmt.Sprintf("Build ended at %v", time.Now()))
	return bContext, nil
}

func (appBuildContext *SingleBuildContext) executeDockerComposeFileIfConfigured(buildContext context.Context) error {
	return nil
}

func (appBuildContext *SingleBuildContext) SaveLogs() error {
	build := mustGetAppBuildFromContext(*appBuildContext.Context)
	build.Logs = appBuildContext.LogsWriter.Logs

	return appBuildContext.Db.Save(&build).Error
}

func (appBuildContext *SingleBuildContext) prepareEnvAndBuildArguments(buildContext context.Context, db *db.AppsDb) error {
	app := mustGetAppFromContext(buildContext)
	os.Setenv("APP_NAME", app.AppName)

	println("Preparing build arguments" + app.BuildVal)
	decoded, err := base64.StdEncoding.DecodeString(app.BuildVal)

	if err != nil {
		return err
	}

	splitted := strings.Split(string(decoded), "\n")

	envs, err := config.ParseEnvLines(splitted)

	if err != nil {
		return err
	}

	config.LoadEnvs(envs)

	envs["APP_NAME"] = app.AppName

	appBuildContext.prepareBuildArgs(envs)

	if app.BuildValMountPath != "" {
		err := os.WriteFile(app.BuildValMountPath, []byte(decoded), 0777)
		if err != nil {
			return err
		}
		appBuildContext.WriteLog(fmt.Sprintf("Build arguments has been written to %s", app.BuildValMountPath))
		log.Printf("Build arguments written to file\n")
	}
	log.Printf("Build arguments prepared\n")

	return err
}

func (appBuildContext *SingleBuildContext) prepareBuildArgs(envs map[string]string) {
	result := ""

	for key := range envs {
		result = result + " --build-arg " + key + " "
	}

	os.Setenv("DOCKER_BUILD_ARGS", result)
}

func (appBuildContext *SingleBuildContext) runPreBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/pre-build.sh", configFolder))

	if err != nil {
		log.Println("Could not read pre-build.sh script.")
		return nil
	}
	return appBuildContext.executeAppCommand(fmt.Sprintf("bash ./%s/pre-build.sh", configFolder))
}

func (appBuildContext *SingleBuildContext) runBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/build.sh", configFolder))

	if err != nil {
		log.Println("Could not read build.sh script.")
		return nil
	}

	return appBuildContext.executeAppCommand(fmt.Sprintf("bash ./%s/build.sh", configFolder))
}

func (appBuildContext *SingleBuildContext) runPostBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/post-build.sh", configFolder))

	if err != nil {
		log.Println("Could not read post-build.sh script.")
		return nil
	}
	return appBuildContext.executeAppCommand(fmt.Sprintf("bash ./%s/post-build.sh", configFolder))
}

type AppLogsWriter struct {
	Logs     string
	LogsChan chan string
}

func (appBuildContext *SingleBuildContext) WriteLog(log string) error {
	_, err := appBuildContext.LogsWriter.Write([]byte(log))
	return err
}

func (w *AppLogsWriter) Write(p []byte) (n int, err error) {
	now := time.Now()

	formattedNow := fmt.Sprintf("%d.%d.%d %d:%d:%d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	newLine := fmt.Sprintf("[%s] %s", formattedNow, string(p))

	w.Logs = w.Logs + newLine
	w.LogsChan <- newLine

	return len(p), nil
}

func (appBuildContext *SingleBuildContext) executeAppCommand(command string) error {
	return util.ExecuteCommand(command, appBuildContext.LogsWriter)
}

func (appBuildContext *SingleBuildContext) cloneRepo(buildContext context.Context) error {
	app := mustGetAppFromContext(buildContext)
	os.Mkdir("/builds", 0777)

	folderName, err := os.MkdirTemp("/builds", "source-*")
	appBuildContext.WriteLog(fmt.Sprintf("created build folder [%s]\n", folderName))

	if err != nil {
		return err
	}

	err = os.Chdir(folderName)

	if err != nil {
		return err
	}

	appBuildContext.WriteLog(fmt.Sprintf("cloning repo %s\n", app.RepoUrl))
	err = appBuildContext.executeAppCommand(fmt.Sprintf("git clone %s .", app.RepoUrl))

	if err != nil {
		return err
	}

	if app.RepoBranch != "" {
		appBuildContext.WriteLog(fmt.Sprintf("switching branch to [%s]\n", app.RepoBranch))
		err = util.ExecuteCommand(fmt.Sprintf("git checkout %s", app.RepoBranch))

		if err != nil {
			return err
		}
	}

	appBuildContext.WriteLog("Build sha: ")
	appBuildContext.executeAppCommand("git rev-parse HEAD")

	return nil
}

func loadBase64EnvFileIntoEnv(buildContext context.Context, db *db.AppsDb) error {
	app := mustGetAppFromContext(buildContext)

	if app.EnvVal != "" {
		err := os.Setenv("BASE_64_ENV_FILE", app.EnvVal)

		if err != nil {
			return err
		}
	}

	return nil
}

func mustGetAppFromContext(ctx context.Context) *types.NanoApp {
	app, ok := ctx.Value(types.CurrentNanoAppContextKey).(*types.NanoApp)

	if !ok {
		panic("could not get app from context")
	}

	return app
}

func mustGetAppBuildFromContext(ctx context.Context) *types.NanoBuild {
	build, ok := ctx.Value(types.CurrentNanoBuildContextKey).(*types.NanoBuild)

	if !ok {
		panic("could not get build from context")
	}

	return build
}
