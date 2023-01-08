package apps

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AppsDb struct {
	*gorm.DB
}

func NewAppsDatabase() (*AppsDb, error) {
	databaseConn, err := gorm.Open(sqlite.Open("apps.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &AppsDb{databaseConn}, nil
}

func (db AppsDb) AutoMigrateModels() error {
	err := db.AutoMigrate(&NanoApp{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&NanoConfig{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&NanoContext{})

	if err != nil {
		return err
	}

	return nil
}

func (db AppsDb) InitConfig() error {
	token := uuid.NewString()

	var any NanoContext
	db.First(&any)

	if any.ID != 0 {
		return nil
	}

	tx := db.Create(&NanoContext{
		Apps: []NanoApp{},
		NanoConfig: NanoConfig{
			GlobalEnvironment: "",
			Token:             token,
		},
	})

	return tx.Error
}