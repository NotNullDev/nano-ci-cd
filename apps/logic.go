package apps

import (
	"cd/config"
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
	buildSecret           = "Qg28syo36sZmnUbpSshKBDbY2wUepp1zXVi5CG6nTyA="
	appSecret             = "shFOAaPW91HkfFd/1SccVGo7aKUV5Zu4MDwEBRgi6pc="
)

func Build(buildArguments BuildArguments) error {
	os.Chdir(baseContainerFolder)

	println(fmt.Sprintf("Build started at %v", time.Now()))

	err := prepareEnvAndBuildArguments(buildArguments)

	err = cloneRepo(buildArguments.RepoUrl)

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

	err = runPostBuildScript()

	if err != nil {
		return err
	}

	err = executeDockerComposeFileIfConfigured(buildArguments.AppName)

	if err != nil {
		return err
	}

	println(fmt.Sprintf("Build ended at %v", time.Now()))
	return nil
}

func executeDockerComposeFileIfConfigured(appName string) error {
	return nil
}

func prepareEnvAndBuildArguments(buildArguments BuildArguments) error {
	os.Setenv("APP_NAME", buildArguments.AppName)
	envs, err := config.ParseEnvFiles(false, "/app/envs/"+buildArguments.AppName)

	if err != nil {
		return err
	}

	config.LoadEnvs(envs)

	envs["APP_NAME"] = buildArguments.AppName
	prepareBuildArgs(envs)
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

func cloneRepo(repoUrl string) error {
	os.Mkdir("builds", 0777)
	folderName, err := os.MkdirTemp("builds", "source-*")

	if err != nil {
		return err
	}

	err = os.Chdir(folderName)

	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", repoUrl, ".")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
