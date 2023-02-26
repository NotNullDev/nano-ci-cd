package api

import (
	"encoding/base64"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	types "github.com/nano-ci-cd/types"
	"github.com/nano-ci-cd/util"
)

func (router Router) ClearBuildFolder(c echo.Context) error {
	err := os.RemoveAll("./build/*")

	if err != nil {
		return err
	}

	err = os.Mkdir("./build", 0755)

	if err != nil {
		return err
	}

	return c.JSON(200, "{}")
}

func (router Router) DownloadDbBackup(c echo.Context) error {
	return c.File("/data/apps.db")
}

func (router Router) ResetGlobalBuildStatus(c echo.Context) error {
	var appConfig types.NanoContext

	if router.db.First(&appConfig).Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: "Context not found",
		})
	}

	appConfig.CurrentlyBuildingAppId = 0

	if router.db.Save(&appConfig).Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: "failed to update context",
		})
	}

	return c.JSON(200, "")
}

func (router Router) DockerSystemPrune(c echo.Context) error {
	err := util.ExecuteCommand("docker system prune -f -a")

	if err != nil {
		return err
	}

	return c.JSON(200, "")
}

func (router Router) UpdateGlobalEnvironment(c echo.Context) error {
	var globalEnv string

	e, err := io.ReadAll(c.Request().Body)

	globalEnv = base64.StdEncoding.EncodeToString(e)

	if err != nil {
		return err
	}

	nanoConfig := types.NanoConfig{}

	tx := router.db.First(&nanoConfig)

	if tx.Error != nil {
		return tx.Error
	}

	nanoConfig.GlobalEnvironment = globalEnv

	router.db.Save(&nanoConfig)

	return c.JSON(200, nanoConfig.GlobalEnvironment)
}
