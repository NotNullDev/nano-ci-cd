package nanocicd

import (
	api "github.com/nano-ci-cd/api"
	db "github.com/nano-ci-cd/db"
)

type NanoCiCD struct {
	Router *api.Router
	Db     *db.AppsDb
}
