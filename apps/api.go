package apps

import (
	"encoding/base64"
	"io"
	"os"
	"strconv"

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

	appCtx.Db.Preload("NanoConfig").Preload("Apps").Find(&apps)

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

type CreateAppRequest struct {
	AppName string `json:"appName"`
}

func (appCtx AppContext) CreateApp(c echo.Context) error {
	var req CreateAppRequest
	c.Bind(&req)

	if req.AppName == "" {
		return c.JSON(400, ErrorResponse{
			Error: "App name is required",
		})
	}

	app := &NanoApp{
		AppName:       req.AppName,
		NanoContextID: 1,
	}

	tx := appCtx.Db.Create(app)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, app)
}

func (appCtx AppContext) UpdateApp(c echo.Context) error {
	var req NanoApp
	err := c.Bind(&req)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	req.NanoContextID = 1

	tx := appCtx.Db.Save(&req)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, req)
}

func (appCtx AppContext) DeleteApp(c echo.Context) error {
	appId := c.Param("id")

	if appId == "" {
		return c.JSON(400, ErrorResponse{
			Error: "App ID is required",
		})
	}

	idAsInt, err := strconv.ParseInt(appId, 10, 64)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	tx := appCtx.Db.Delete(&NanoApp{}, idAsInt)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, idAsInt)
}
