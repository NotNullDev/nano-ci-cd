package main

import (
	"cd/config"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	configFolder = ".nano-cicd"
)

var (
	globalEnv = make(map[string]string)
)

func init() {
	env, err := config.ParseEnvFiles(false, "envs/.global")

	if err != nil {
		panic(err.Error())
	}

	globalEnv = env

	config.LoadEnvs(env)
}

type BuildArguments struct {
	RepoUrl string `json:"repoUrl"`
	AppName string `json:"appName"`
	// DockerfilePath           string   `json:"dockerfilePath"`
	// PreDeployScriptLocation  string   `json:"preDeployScriptLocation"`
	// DeployScriptLocation     string   `json:"deployScriptLocation"`
	// PostDeployScriptLocation string   `json:"postDeployScriptLocation"`
	// EnvFileNames             []string `json:"envFileNames"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	e := echo.New()

	e.POST("/build", build)

	if err := e.Start(":8080"); err != nil {
		panic(err.Error())
	}
}

func build(c echo.Context) error {
	var buildArguments BuildArguments

	err := c.Bind(&buildArguments)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	envs, err := config.ParseEnvFiles(false, "envs/"+buildArguments.AppName)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}
	config.LoadEnvs(envs)

	err = cloneRepo(buildArguments.RepoUrl)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	err = runPreBuildScript()

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	err = runBuildScript()

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	err = runPostBuildScript()

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(200, "")
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
