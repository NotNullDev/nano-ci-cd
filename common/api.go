package common

import "github.com/labstack/echo/v4"

type AppContext struct {
	Echo *echo.Echo
	Db   *AppsDb
}
