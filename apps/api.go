package apps

import (
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AppContext struct {
	Echo *echo.Echo
	Db   *AppsDb
}

type AppsApi struct {
	*echo.Echo
}

func (appCtx AppContext) HandlePostRequest(c echo.Context) error {
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

func (appCtx AppContext) GetNanoContext(c echo.Context) error {
	var apps NanoContext

	appCtx.Db.Preload("NanoConfig").Find(&apps)

	return c.JSON(200, apps)
}

func (appCtx AppContext) ResetToken(c echo.Context) error {
	token := uuid.NewString()

	nanoConfig := NanoConfig{}

	tx := appCtx.Db.First(&nanoConfig)

	if tx.Error != nil {
		return tx.Error
	}

	nanoConfig.Token = token

	appCtx.Db.Save(&nanoConfig)

	return c.JSON(200, token)
}

func (appCtx AppContext) UpdateGlobalEnvironment(c echo.Context) error {
	var globalEnv string

	e, err := io.ReadAll(c.Request().Body)

	globalEnv = base64.StdEncoding.EncodeToString(e)

	if err != nil {
		return err
	}

	nanoConfig := NanoConfig{}

	tx := appCtx.Db.First(&nanoConfig)

	if tx.Error != nil {
		return tx.Error
	}

	nanoConfig.GlobalEnvironment = globalEnv

	appCtx.Db.Save(&nanoConfig)

	return c.JSON(200, nanoConfig.GlobalEnvironment)
}

func (c *AppContext) createApp(appName string) error {
	if appName == "" {
		return errors.New("appName is required")
	}

	app := NanoApp{
		AppName: appName,
	}

	tx := c.Db.Create(&app)

	return tx.Error
}
