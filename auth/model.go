package auth

import (
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
