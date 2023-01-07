package main

import (
	"cd/apps"
	"cd/config"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	latestShaShortCommand = "git rev-parse --short HEAD"
	baseContainerFolder   = "/app"
)

var (
	globalEnv = make(map[string]string)
)

func init() {
	env, err := config.ParseEnvFiles(false, "/app/envs/.global")

	if err != nil {
		panic(err.Error())
	}

	globalEnv = env

	config.LoadEnvs(env)
}

func main() {
	db, err := apps.NewAppsDatabase()

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrateModels()

	err = db.InitConfig()

	if err != nil {
		panic(err.Error())
	}

	e := echo.New()

	app := apps.AppContext{
		Echo: e,
		Db:   db,
	}

	// dashboard
	e.GET("/", app.GetNanoContext)
	e.POST("/reset-token", app.ResetToken)
	e.POST("/update-global-env", app.UpdateGlobalEnvironment)
	e.POST("/create-app", app.CreateApp)
	e.POST("/update-app", app.UpdateApp)
	e.POST("/delete-app", app.UpdateApp)

	// build trigger
	e.POST("/build", app.HandlePostRequest)

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		println(err.Error() + " at " + time.Now().String())
		ctx.JSON(500, err.Error())
	}

	if err := e.Start(":8080"); err != nil {
		panic(err.Error())
	}
}
