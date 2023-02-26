package main

import (
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	router "github.com/nano-ci-cd/api"
	"github.com/nano-ci-cd/auth"
	"github.com/nano-ci-cd/config"
	db "github.com/nano-ci-cd/db"
	"github.com/nano-ci-cd/nanocicd"
	"github.com/nano-ci-cd/types"
)

func init() {
	env, err := config.ParseEnvFiles(false, "/app/envs/.global")

	if err != nil {
		panic(err.Error())
	}

	config.LoadEnvs(env)
}

func main() {
	db, err := db.NewAppsDatabase()

	if err != nil {
		panic(err.Error())
	}

	prepareDatabase(db)

	e := echo.New()
	initMiddleware(e, db)

	apiRouter := router.NewRouter(e, db)

	app := nanocicd.NanoCiCD{
		Router: apiRouter,
		Db:     db,
	}

	nanoContext := types.NanoContext{}

	db.First(&nanoContext)
	nanoContext.CurrentlyBuildingAppId = 0
	db.Save(&nanoContext)

	// dashboard
	e.GET("/", apiRouter.GetNanoContext)
	e.POST("/reset-token", apiRouter.ResetToken)

	e.POST("/create-app", apiRouter.CreateApp)
	e.POST("/update-app", apiRouter.UpdateApp)
	e.POST("/login", apiRouter.Login)
	e.POST("/update-user", apiRouter.UpdateUser)
	e.DELETE("/delete-app", apiRouter.DeleteApp)
	e.GET("/available-builds-metadata", apiRouter.GetBuilds)
	e.GET("/build", apiRouter.GetBuild)
	e.GET("/logs", apiRouter.GetLogs)

	e.GET("/clear-builds", apiRouter.ClearBuildFolder)
	e.GET("/download-backup", apiRouter.DownloadDbBackup)
	e.GET("/reset-global-build-status", apiRouter.ResetGlobalBuildStatus)
	e.GET("/docker-system-prune", apiRouter.DockerSystemPrune)
	e.POST("/update-global-env", apiRouter.UpdateGlobalEnvironment)

	// build trigger
	e.POST("/build", apiRouter.HandlePostRequest)

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		println(err.Error() + " at " + time.Now().String())
		ctx.JSON(500, err.Error())
	}

	if err := app.Router.Echo.Start(":8080"); err != nil {
		panic(err.Error())
	}
}

func prepareDatabase(db *db.AppsDb) {
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

func initMiddleware(e *echo.Echo, db *db.AppsDb) {
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/login" {
				return next(c)
			}

			if c.Path() == "/build" && strings.ToLower(c.Request().Method) == "post" {
				var config types.NanoConfig
				db.First(&config)

				h := c.Request().Header
				key := h.Get("Authorization")

				if key != config.Token {
					log.Printf("invalid token: %s", key)
					return c.JSON(403, types.ErrorResponse{
						Error: "invalid token",
					})
				}

				return next(c)
			}

			token := c.Request().Header.Get("nano-token")

			if token == "" {
				return c.JSON(401, types.ErrorResponse{
					Error: "missing token",
				})
			}

			err := auth.ValidateToken(token)

			if err != nil {
				return c.JSON(401, types.ErrorResponse{
					Error: "invalid token",
				})
			}

			session := &types.NanoSession{
				Token: token,
			}

			tx := db.DB.First(&session)

			if tx.Error != nil {
				return c.JSON(401, types.ErrorResponse{
					Error: "invalid token",
				})
			}

			if session.ID == 0 {
				return c.JSON(403, types.ErrorResponse{
					Error: "session not found",
				})
			}

			return next(c)
		}
	})
}
