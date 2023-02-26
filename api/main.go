package api

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nano-ci-cd/auth"
	"github.com/nano-ci-cd/config"
	db "github.com/nano-ci-cd/db"
	logic "github.com/nano-ci-cd/logic"
	types "github.com/nano-ci-cd/types"
	"golang.org/x/crypto/bcrypt"
)

type Router struct {
	Echo *echo.Echo
	db   *db.AppsDb
}

func NewRouter(e *echo.Echo, db *db.AppsDb) *Router {
	return &Router{
		Echo: e,
		db:   db,
	}
}

func (appCtx Router) HandlePostRequest(c echo.Context) error {
	// setSSEHeaders(c)

	nanoContext := types.NanoContext{}
	appCtx.db.First(&nanoContext)

	if nanoContext.CurrentlyBuildingAppId != 0 {
		return c.JSON(400, types.ErrorResponse{
			Error: "another app is currently building",
		})
	}

	var appConfig types.NanoConfig

	tx := appCtx.db.First(&appConfig)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	if appConfig.Token != c.Request().Header.Get("Authorization") {
		return c.JSON(403, types.ErrorResponse{
			Error: "invalid token",
		})
	}

	appName := c.QueryParam("appName")

	if appName == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "missing appName",
		})
	}

	app := types.NanoApp{
		AppName: appName,
	}

	tx = appCtx.db.First(&app, "app_name = ?", appName)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	print("Found app name " + app.AppName)

	if app.AppStatus != "enabled" {
		return c.JSON(400, types.ErrorResponse{
			Error: "app is disabled",
		})
	}

	c.JSON(200, "")
	buildContext := context.Background()
	buildContext = context.WithValue(buildContext, types.CurrentNanoAppContextKey, &app)

	err := loadGlobalEnvs(appConfig)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}
	logsChan := make(chan string)
	done := make(chan bool)

	savedbuff := []string{}
	go func() {
		build := &types.NanoBuild{
			AppID:       app.ID,
			StartedAt:   time.Now(),
			BuildStatus: "running",
		}
		appCtx.db.Create(&build)

		buildContext = context.WithValue(buildContext, types.CurrentNanoBuildContextKey, build)

		{
			appCtx.db.First(&nanoContext)
			nanoContext.CurrentlyBuildingAppId = app.ID
			appCtx.db.Save(&nanoContext)
		}

		defer func() {
			appCtx.db.First(&nanoContext)
			nanoContext.CurrentlyBuildingAppId = 0
			appCtx.db.Save(&nanoContext)
		}()

		go func() {
		outer:
			for {
				select {
				case log := <-logsChan:
					// c.Response().Write([]byte(log))
					// savedbuff = append(savedbuff, log)
					build.Logs = build.Logs + log
					appCtx.db.Save(&build)
					// os.Stderr.Write([]byte("haha: " + log + "\n"))
				case <-done:
					break outer
				}
			}
		}()

		buildContext, err := logic.Build(buildContext, appCtx.db, logsChan)

		if err != nil {
			println(err.Error())
			build.BuildStatus = "failed"
			return
		} else {
			build.BuildStatus = "success"
		}

		appCtx.db.Save(&build)

		buildContext.SaveLogs()
		done <- true
		println("haha saved buf: ", savedbuff)
	}()

	return c.JSON(200, "")
}

func (appCtx Router) GetNanoContext(c echo.Context) error {
	var apps types.NanoContext

	appCtx.db.Preload("NanoConfig").Preload("Apps").Find(&apps)

	return c.JSON(200, apps)
}

func (appCtx Router) ResetToken(c echo.Context) error {
	token := uuid.NewString()

	nanoConfig := types.NanoConfig{}

	tx := appCtx.db.First(&nanoConfig)

	if tx.Error != nil {
		return tx.Error
	}

	nanoConfig.Token = token

	appCtx.db.Save(&nanoConfig)

	return c.JSON(200, token)
}

type CreateAppRequest struct {
	AppName string `json:"appName"`
}

func (appCtx Router) CreateApp(c echo.Context) error {
	var req CreateAppRequest
	c.Bind(&req)

	if req.AppName == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "App name is required",
		})
	}

	app := &types.NanoApp{
		AppName:       req.AppName,
		NanoContextID: 1,
		AppStatus:     "enabled",
	}

	tx := appCtx.db.Create(app)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, app)
}

func (appCtx Router) UpdateApp(c echo.Context) error {
	var req types.NanoApp
	err := c.Bind(&req)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	req.NanoContextID = 1
	req.BuildVal = base64.StdEncoding.EncodeToString([]byte(req.BuildVal))
	req.EnvVal = base64.StdEncoding.EncodeToString([]byte(req.EnvVal))

	tx := appCtx.db.Save(&req)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, req)
}

func (appCtx Router) DeleteApp(c echo.Context) error {
	appId := c.QueryParam("id")

	if appId == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "App ID is required",
		})
	}

	idAsInt, err := strconv.ParseInt(appId, 10, 64)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	tx := appCtx.db.Delete(&types.NanoApp{}, idAsInt)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, idAsInt)
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (appCtx Router) UpdateUser(c echo.Context) error {
	token := c.Request().Header.Get("nano-token")

	var req UpdateUserRequest

	err := c.Bind(&req)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	session := &types.NanoSession{
		Token: token,
	}

	appCtx.db.First(&session)

	if session.ID == 0 {
		return c.JSON(403, types.ErrorResponse{
			Error: "Invalid token",
		})
	}

	user := &types.NanoUser{}

	appCtx.db.First(&user, session.NanoUserID)

	if user.ID == 0 {
		return c.JSON(403, types.ErrorResponse{
			Error: "Invalid token",
		})
	}

	user.Username = req.Username
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	user.Password = string(hashed)

	tx := appCtx.db.Save(&user)

	if tx.Error != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, "")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (appCtx Router) Login(c echo.Context) error {
	var req LoginRequest
	c.Bind(&req)

	if req.Username == "" || req.Password == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "Username or password is empty",
		})
	}

	user := types.NanoUser{
		Username: req.Username,
	}

	appCtx.db.First(&user)

	if user.ID == 0 {
		return c.JSON(403, types.ErrorResponse{
			Error: "Invalid username or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return c.JSON(403, types.ErrorResponse{
			Error: "Invalid username or password",
		})
	}

	token, err := auth.CreateToken()

	if err != nil {
		return c.JSON(500, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	session := &types.NanoSession{
		Token:      token,
		NanoUserID: user.ID,
	}

	tx := appCtx.db.Save(&session)

	if tx.Error != nil {
		return c.JSON(500, types.ErrorResponse{
			Error: "",
		})
	}

	return c.JSON(200, token)
}

func (appCtx Router) GetLogs(c echo.Context) error {
	appId := c.QueryParam("appId")

	if appId == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "App ID is required",
		})
	}

	idAsInt, err := strconv.ParseInt(appId, 10, 64)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	logs := &types.NanoBuild{}

	appCtx.db.Order("id desc").Where("app_id = ?", idAsInt).Find(&logs).Limit(1)

	if logs.ID == 0 {
		return c.JSON(400, types.ErrorResponse{
			Error: "Logs for the requested app not found",
		})
	}

	return c.JSON(200, logs)
}

type GetBuildsResponseEntity struct {
	ID        int64  `json:"id"`
	StartedAt string `json:"date" gorm:"started_at"`
}

func (appCtx Router) GetBuilds(c echo.Context) error {
	builds := []GetBuildsResponseEntity{}

	appCtx.db.Raw("select n.id, n.started_at from nano_builds n order by n.started_at desc").Scan(&builds)

	return c.JSON(200, builds)
}

func (appCtx Router) GetBuild(c echo.Context) error {
	appId := c.QueryParam("buildId")

	if appId == "" {
		return c.JSON(400, types.ErrorResponse{
			Error: "App ID is required",
		})
	}

	idAsInt, err := strconv.ParseInt(appId, 10, 64)

	if err != nil {
		return c.JSON(400, types.ErrorResponse{
			Error: err.Error(),
		})
	}

	build := types.NanoBuild{}

	appCtx.db.Where("id = ?", idAsInt).First(&build)

	return c.JSON(200, build)
}

func loadGlobalEnvs(appConfig types.NanoConfig) error {
	decoded, err := base64.StdEncoding.DecodeString(appConfig.GlobalEnvironment)

	if err != nil {
		return err
	}

	splitted := strings.Split(string(decoded), "\n")

	envs, err := config.ParseEnvLines(splitted)

	if err != nil {
		return err
	}

	config.LoadEnvs(envs)
	return nil
}

func setSSEHeaders(c echo.Context) {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
}
