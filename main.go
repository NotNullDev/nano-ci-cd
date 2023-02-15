package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nano-ci-cd/apps"
	"github.com/nano-ci-cd/auth"
	"github.com/nano-ci-cd/config"
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

	prepareDatabase(db)

	e := echo.New()
	initMiddleware(e, db)

	app := apps.AppContext{
		Echo: e,
		Db:   db,
	}

	nanoContext := apps.NanoContext{}

	db.First(&nanoContext)
	nanoContext.CurrentlyBuildingAppId = 0
	db.Save(&nanoContext)

	// dashboard
	e.GET("/", app.GetNanoContext)
	e.POST("/reset-token", app.ResetToken)
	e.POST("/update-global-env", app.UpdateGlobalEnvironment)
	e.POST("/create-app", app.CreateApp)
	e.POST("/update-app", app.UpdateApp)
	e.DELETE("/delete-app", app.DeleteApp)
	e.GET("/clear-builds", app.ClearBuildFolder)
	e.GET("/download-backup", app.DownloadDbBackup)
	e.POST("/login", app.Login)
	e.POST("/update-user", app.UpdateUser)
	e.GET("/logs", app.GetLogs)
	e.GET("/reset-global-build-status", app.ResetGlobalBuildStatus)
	e.GET("/docker-system-prune", app.DockerSystemPrune)

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

func prepareDatabase(db *apps.AppsDb) {
	err := db.AutoMigrateModels()

	if err != nil {
		panic(err.Error())
	}

	err = db.InitConfig()

	if err != nil {
		panic(err.Error())
	}

	err = db.InitUser()

	if err != nil {
		panic(err.Error())
	}
}

func initMiddleware(e *echo.Echo, db *apps.AppsDb) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.Secure())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/login" {
				return next(c)
			}

			if c.Path() == "/build" {
				var config apps.NanoConfig
				db.First(&config)

				h := c.Request().Header
				key := h.Get("Authorization")

				if key != config.Token {
					log.Printf("invalid token: %s", key)
					return c.JSON(403, apps.ErrorResponse{
						Error: "invalid token",
					})
				}

				return next(c)
			}

			token := c.Request().Header.Get("nano-token")

			if token == "" {
				return c.JSON(401, apps.ErrorResponse{
					Error: "missing token",
				})
			}

			err := auth.ValidateToken(token)

			if err != nil {
				return c.JSON(401, apps.ErrorResponse{
					Error: "invalid token",
				})
			}

			session := &auth.NanoSession{
				Token: token,
			}

			tx := db.DB.First(&session)

			if tx.Error != nil {
				return c.JSON(401, apps.ErrorResponse{
					Error: "invalid token",
				})
			}

			if session.ID == 0 {
				return c.JSON(403, apps.ErrorResponse{
					Error: "session not found",
				})
			}

			return next(c)
		}
	})
}
