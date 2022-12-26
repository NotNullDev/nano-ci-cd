package main

import (
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

type BuildArguments struct {
	RepoUrl string `json:"repoUrl"`
	// DockerfilePath           string   `json:"dockerfilePath"`
	// PreDeployScriptLocation  string   `json:"preDeployScriptLocation"`
	// DeployScriptLocation     string   `json:"deployScriptLocation"`
	// PostDeployScriptLocation string   `json:"postDeployScriptLocation"`
	// EnvFileNames             []string `json:"envFileNames"`
	// AppName                  string   `json:"appName"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	e := echo.New()

	e.POST("/build", build)

	e.Start(":8080")
}

func build(c echo.Context) error {
	var buildArguments BuildArguments

	err := c.Bind(&buildArguments)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	err = cloneRepo(buildArguments.RepoUrl)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	err = runPreDeployScript()

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	runDeployScript()

	err = runPostDeployScript()

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(200, "")
}

func runPreDeployScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/pre-deploy.sh", configFolder))

	if err != nil {
		log.Println("Could not read pre-deploy.sh script.")
		return nil
	}
	return executeCommand(fmt.Sprintf("bash ./%s/pre-deploy.sh", configFolder))
}

func runDeployScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/deploy.sh", configFolder))

	if err != nil {
		log.Println("Could not read deploy.sh script.")
		return nil
	}

	return executeCommand(fmt.Sprintf("bash ./%s/deploy.sh", configFolder))
}

func runPostDeployScript() error {
	_, err := os.Stat(fmt.Sprintf("./%s/post-deploy.sh", configFolder))

	if err != nil {
		log.Println("Could not read post-deploy.sh script.")
		return nil
	}
	return executeCommand(fmt.Sprintf("bash ./%s/post-deploy.sh", configFolder))
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
