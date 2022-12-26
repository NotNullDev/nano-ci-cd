package main

import (
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
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

	content, err := os.ReadFile("Dockerfile")
	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	println(string(content))

	return c.JSON(200, "")
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
