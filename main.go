package main

import (
	"cd/apps"
	"cd/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	env, err := config.ParseEnvFiles(false, "/app/envs/.global")

	if err != nil {
		panic(err.Error())
	}

	config.LoadEnvs(env)
}

func main() {
	db, err := apps.NewAppsDatabase()

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrateModels()

	if err != nil {
		panic(err.Error())
	}

	err = db.InitConfig()

	if err != nil {
		panic(err.Error())
	}

	e := echo.New()

	e.Use(middleware.CORS())

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
	e.DELETE("/delete-app", app.DeleteApp)

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
