package db

import (
	"fmt"
	"time"

	config "github.com/kurdilesmana/go-saving-api/apps/account/configs"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPgsqlConnection(dbConfig *config.ACCDatabase, logger *logging.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.ACCHost,
		dbConfig.ACCUsername,
		dbConfig.ACCPassword,
		dbConfig.ACCDBName,
		dbConfig.ACCPort,
	)

	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to access database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.ACCMaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.ACCMaxConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(dbConfig.ACCConnMaxLifetime))

	// Auto Migrate the Account struct
	log.Info(logrus.Fields{}, nil, "start migrate database...")
	db.AutoMigrate(&accountModel.Account{})

	return db, nil
}
