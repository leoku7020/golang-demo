package mysqlkit

import (
	dsn "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"demo/pkg/logger"
)

func NewConfig(dsnStr string) (*dsn.Config, error) {
	return dsn.ParseDSN(dsnStr)
}

func NewGORM(dsn string, plugins ...gorm.Plugin) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                   logger.NewGormLogger(),
		DisableNestedTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins {
		if err := db.Use(plugin); err != nil {
			return nil, err
		}
	}

	return db, nil
}
