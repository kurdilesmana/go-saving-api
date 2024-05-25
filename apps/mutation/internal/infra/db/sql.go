package db

import (
	"fmt"
	"time"

	config "github.com/kurdilesmana/go-saving-api/apps/mutation/configs"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPgsqlConnection(dbConfig *config.MUTDatabase, logger *logging.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.MUTHost,
		dbConfig.MUTUsername,
		dbConfig.MUTPassword,
		dbConfig.MUTDBName,
		dbConfig.MUTPort,
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

	sqlDB.SetMaxIdleConns(dbConfig.MUTMaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.MUTMaxConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(dbConfig.MUTConnMaxLifetime))

	// Auto Migrate the Account struct
	log.Info(logrus.Fields{}, nil, "start migrate database...")
	db.AutoMigrate(&transactionModel.Transaction{})
	db.AutoMigrate(&transactionModel.TransactionDetail{})

	return db, nil
}
