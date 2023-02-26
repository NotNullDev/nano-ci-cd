package types

import (
	"time"

	"gorm.io/gorm"
)

type NanoUser struct {
	gorm.Model `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type NanoSession struct {
	gorm.Model
	NanoUserID uint
	Token      string
}

type NanoSessionData struct {
	gorm.Model
	NanoSessionID uint
	Key           string
	Value         string
}

type NanoBuild struct {
	gorm.Model
	AppID       uint      `json:"appId"`
	BuildStatus string    `json:"buildStatus"` // running, failed, success
	Logs        string    `json:"logs"`
	StartedAt   time.Time `json:"startedAt"`
	FinishedAt  time.Time `json:"finishedAt"`
}
