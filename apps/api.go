package apps

import (
	"errors"
	"os"

	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	Echo echo.Context
	Db   AppsDb
}

type AppsApi struct {
	*echo.Echo
}

func HandlePostRequest(c echo.Context) error {
	var buildArguments BuildArguments

	if dev := os.Getenv("DEV_MODE"); dev != "" {
		buildArguments = BuildArguments{
			RepoUrl: "https://gitea.notnulldev.com/notnulldev/nano-ci-cd",
			AppName: "nano-ci-cd",
		}
	} else {
		var giteaArgs GiteaHook

		err := c.Bind(&giteaArgs)

		if err != nil {
			return c.JSON(400, ErrorResponse{
				Error: err.Error(),
			})
		}

		buildArguments = BuildArguments{
			RepoUrl: giteaArgs.Repository.CloneURL,
			AppName: giteaArgs.Repository.Name,
		}
	}

	err := Build(buildArguments)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(200, "")
}

func (c *CustomContext) createApp(appName string) error {
	if appName == "" {
		return errors.New("appName is required")
	}

	app := NanoApp{
		AppName: appName,
	}

	tx := c.db.Create(&app)

	return tx.Error
}
