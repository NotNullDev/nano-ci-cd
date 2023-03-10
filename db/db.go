package apps

import (
	"os"

	"github.com/glebarez/sqlite"
	"github.com/nano-ci-cd/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppsDb struct {
	*gorm.DB
}

func NewAppsDatabase() (*AppsDb, error) {
	os.Mkdir("/data", 0777)
	databaseConn, err := gorm.Open(sqlite.Open("/data/apps.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &AppsDb{databaseConn}, nil
}

func (db AppsDb) AutoMigrateModels() error {
	err := db.AutoMigrate(&types.NanoApp{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoConfig{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoContext{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoUser{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoSession{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoSessionData{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.NanoBuild{})

	if err != nil {
		return err
	}

	return nil
}

func (db AppsDb) InitConfig() error {
	token := "62285a21-547d-46db-a9fd-a2fec5161da5" // hardcoded initial token on both server and client

	var any types.NanoContext
	db.First(&any)

	if any.ID != 0 {
		return nil
	}

	initToken := os.Getenv("NANO_INIT_TOKEN")

	if initToken != "" {
		token = initToken
	}

	tx := db.Create(&types.NanoContext{
		Apps: []types.NanoApp{},
		NanoConfig: types.NanoConfig{
			GlobalEnvironment: "",
			Token:             token,
		},
	})

	return tx.Error
}

func (db AppsDb) InitUser() error {
	var any types.NanoUser
	db.First(&any)

	if any.ID != 0 {
		return nil
	}

	username := "admin"
	password := "admin"

	if os.Getenv("NANO_INIT_USERNAME") != "" {
		username = os.Getenv("NANO_INIT_USERNAME")
	}

	if os.Getenv("NANO_INIT_PASSWORD") != "" {
		password = os.Getenv("NANO_INIT_PASSWORD")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	tx := db.Create(&types.NanoUser{
		Username: username,
		Password: string(hash),
	})

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
