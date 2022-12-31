package main

import (
	"cd/apps"
	"cd/config"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	configFolder          = ".nano-cicd"
	latestShaShortCommand = "git rev-parse --short HEAD"
	baseContainerFolder   = "/app"
	buildSecret           = "Qg28syo36sZmnUbpSshKBDbY2wUepp1zXVi5CG6nTyA="
	appSecret             = "shFOAaPW91HkfFd/1SccVGo7aKUV5Zu4MDwEBRgi6pc="
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

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &apps.CustomContext{Echo: c, Db: db}
			return next(cc)
		}
	})

	e.POST("/build", apps.HandlePostRequest)

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		println(err.Error() + " at " + time.Now().String())
	}

	if err := e.Start(":8080"); err != nil {
		panic(err.Error())
	}
}
