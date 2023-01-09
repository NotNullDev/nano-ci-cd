package apps

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nano-ci-cd/auth"
	"github.com/nano-ci-cd/config"
	"golang.org/x/crypto/bcrypt"
)

type AppContext struct {
	Echo *echo.Echo
	Db   *AppsDb
}

type AppsApi struct {
	*echo.Echo
}

type ContextKey string

var contextKey ContextKey = "app"

func (appCtx AppContext) HandlePostRequest(c echo.Context) error {
	var appConfig NanoConfig

	tx := appCtx.Db.First(&appConfig)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	if appConfig.Token != c.Request().Header.Get("Authorization") {
		return c.JSON(403, ErrorResponse{
			Error: "invalid token",
		})
	}

	appName := c.QueryParam("appName")

	if appName == "" {
		return c.JSON(400, ErrorResponse{
			Error: "missing appName",
		})
	}

	app := NanoApp{
		AppName: appName,
	}

	tx = appCtx.Db.First(&app, "app_name = ?", appName)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	print("Found app name " + app.AppName)

	if app.AppStatus != "enabled" {
		return c.JSON(400, ErrorResponse{
			Error: "app is disabled",
		})
	}

	c.JSON(200, "")
	buildContext := context.Background()
	buildContext = context.WithValue(buildContext, contextKey, &app)

	err := loadGlobalEnvs(appConfig)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}
	err = Build(buildContext, appCtx.Db)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	return nil
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
		AppStatus:     "enabled",
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
	req.BuildVal = base64.StdEncoding.EncodeToString([]byte(req.BuildVal))
	req.EnvVal = base64.StdEncoding.EncodeToString([]byte(req.EnvVal))

	tx := appCtx.Db.Save(&req)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, req)
}

func (appCtx AppContext) DeleteApp(c echo.Context) error {
	appId := c.QueryParam("id")

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

func (appCtx AppContext) ClearBuildFolder(c echo.Context) error {
	err := os.RemoveAll("./build")

	if err != nil {
		return err
	}

	err = os.Mkdir("./build", 0755)

	if err != nil {
		return err
	}

	return c.JSON(200, "{}")
}

func (appCtx AppContext) DownloadDbBackup(c echo.Context) error {
	return c.File("./apps.db")
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (appCtx AppContext) UpdateUser(c echo.Context) error {
	token := c.Request().Header.Get("nano-token")

	var req UpdateUserRequest

	err := c.Bind(&req)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	session := &auth.NanoSession{
		Token: token,
	}

	appCtx.Db.First(&session)

	if session.ID == 0 {
		return c.JSON(403, ErrorResponse{
			Error: "Invalid token",
		})
	}

	user := &auth.NanoUser{}

	appCtx.Db.First(&user, session.NanoUserID)

	if user.ID == 0 {
		return c.JSON(403, ErrorResponse{
			Error: "Invalid token",
		})
	}

	user.Username = req.Username
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.JSON(400, ErrorResponse{
			Error: err.Error(),
		})
	}

	user.Password = string(hashed)

	tx := appCtx.Db.Save(&user)

	if tx.Error != nil {
		return c.JSON(400, ErrorResponse{
			Error: tx.Error.Error(),
		})
	}

	return c.JSON(200, "")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (appCtx AppContext) Login(c echo.Context) error {
	var req LoginRequest
	c.Bind(&req)

	if req.Username == "" || req.Password == "" {
		return c.JSON(400, ErrorResponse{
			Error: "Username or password is empty",
		})
	}

	user := auth.NanoUser{
		Username: req.Username,
	}

	appCtx.Db.First(&user)

	if user.ID == 0 {
		return c.JSON(403, ErrorResponse{
			Error: "Invalid username or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return c.JSON(403, ErrorResponse{
			Error: "Invalid username or password",
		})
	}

	token, err := auth.CreateToken()

	if err != nil {
		return c.JSON(500, ErrorResponse{
			Error: err.Error(),
		})
	}

	session := &auth.NanoSession{
		Token:      token,
		NanoUserID: user.ID,
	}

	tx := appCtx.Db.Save(&session)

	if tx.Error != nil {
		return c.JSON(500, ErrorResponse{
			Error: "",
		})
	}

	return c.JSON(200, token)
}

func loadGlobalEnvs(appConfig NanoConfig) error {
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
