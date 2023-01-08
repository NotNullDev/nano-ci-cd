package apps

import (
	"cd/config"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	configFolder          = ".nano-cicd"
	latestShaShortCommand = "git rev-parse --short HEAD"
	baseContainerFolder   = "/app"
)

func Build(buildContext context.Context, db *AppsDb) error {
	os.Chdir(baseContainerFolder)

	println(fmt.Sprintf("Build started at %v", time.Now()))

	err := cloneRepo(buildContext)

	if err != nil {
		return err
	}

	err = prepareEnvAndBuildArguments(buildContext, db)

	if err != nil {
		return err
	}

	err = runPreBuildScript()

	if err != nil {
		return err
	}

	err = runBuildScript()

	if err != nil {
		return err
	}

	err = loadBase64EnvFileIntoEnv(buildContext, db)

	if err != nil {
		return err
	}

	err = runPostBuildScript()

	if err != nil {
		return err
	}

	err = executeDockerComposeFileIfConfigured(buildContext)

	if err != nil {
		return err
	}

	println(fmt.Sprintf("Build ended at %v", time.Now()))
	return nil
}

func executeDockerComposeFileIfConfigured(buildContext context.Context) error {
	return nil
}

func prepareEnvAndBuildArguments(buildContext context.Context, db *AppsDb) error {
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
	prepareBuildArgs(envs)

	if app.BuildValMountPath != "" {
		err := os.WriteFile(app.BuildValMountPath, []byte(decoded), 0777)
		if err != nil {
			return err
		}
	}

	return err
}

func prepareBuildArgs(envs map[string]string) {
	result := ""

	for key := range envs {
		result = result + " --build-arg " + key + " "
	}

	os.Setenv("DOCKER_BUILD_ARGS", result)
}

func runPreBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/pre-build.sh", configFolder))

	if err != nil {
		log.Println("Could not read pre-build.sh script.")
		return nil
	}
	return executeCommand(fmt.Sprintf("bash ./%s/pre-build.sh", configFolder))
}

func runBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/build.sh", configFolder))

	if err != nil {
		log.Println("Could not read build.sh script.")
		return nil
	}

	return executeCommand(fmt.Sprintf("bash ./%s/build.sh", configFolder))
}

func runPostBuildScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/post-build.sh", configFolder))

	if err != nil {
		log.Println("Could not read post-build.sh script.")
		return nil
	}
	return executeCommand(fmt.Sprintf("bash ./%s/post-build.sh", configFolder))
}

func executeCommand(command string) error {
	splitted := strings.Split(command, " ")

	if len(splitted) <= 1 {
		return errors.New("could not split command")
	}

	cmd := exec.Command(splitted[0], splitted[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func cloneRepo(buildContext context.Context) error {
	app := mustGetAppFromContext(buildContext)
	os.Mkdir("builds", 0777)

	folderName, err := os.MkdirTemp("builds", "source-*")

	if err != nil {
		return err
	}

	err = os.Chdir(folderName)

	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", app.RepoUrl, ".")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	if app.RepoBranch != "" {
		err = executeCommand(fmt.Sprintf("git checkout %s", app.RepoBranch))

		if err != nil {
			return err
		}
	}

	return nil
}

func loadBase64EnvFileIntoEnv(buildContext context.Context, db *AppsDb) error {
	app := mustGetAppFromContext(buildContext)

	if app.EnvVal != "" {
		err := os.Setenv("BASE_64_ENV_FILE", app.EnvVal)

		if err != nil {
			return err
		}
	}

	return nil
}

func mustGetAppFromContext(ctx context.Context) *NanoApp {
	app, ok := ctx.Value(contextKey).(*NanoApp)

	if !ok {
		panic("could not get app from context")
	}

	return app
}
