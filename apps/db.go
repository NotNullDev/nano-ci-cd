package apps

import (
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

	return &AppsDb{DB: databaseConn}, nil
}

func (db AppsDb) AutoMigrateModels() error {
	err := db.AutoMigrate(&NanoApp{})

	if err != nil {
		return err
	}

	return nil
}
