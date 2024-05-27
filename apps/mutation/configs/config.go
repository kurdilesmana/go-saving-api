package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/spf13/cast"
)

type EnvironmentConfig struct {
	Env         string
	AppConfig   AppConfig
	MUTDatabase MUTDatabase
	Cache       Redis
	Log         *logging.Logger
}

type AppConfig struct {
	Name           string
	Version        string
	Port           int
	MaxRequestTime int
}

type Database struct {
	*sqlx.DB
}

type Transaction struct {
	*sqlx.Tx
}

func LoadENVConfig() (config EnvironmentConfig, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		err = fmt.Errorf("failed to load .env file: %w", err)
		return EnvironmentConfig{}, err
	}

	port, err := strconv.Atoi(os.Getenv("MUT_APP_PORT"))
	if err != nil {
		err = fmt.Errorf("error when convert string to int: %w", err)
	}

	config = EnvironmentConfig{
		Env: os.Getenv("ENV"),
		AppConfig: AppConfig{
			Name:           os.Getenv("MUT_APP_NAME"),
			Version:        os.Getenv("MUT_APP_VERSION"),
			Port:           port,
			MaxRequestTime: cast.ToInt(os.Getenv("MUT_APP_MAX_REQUEST_TIME")),
		},
		MUTDatabase: MUTDatabase{
			MUTEngine:          os.Getenv("MUT_DATABASE_ENGINE"),
			MUTHost:            os.Getenv("MUT_DATABASE_HOST"),
			MUTPort:            cast.ToInt(os.Getenv("MUT_DATABASE_PORT")),
			MUTUsername:        os.Getenv("MUT_DATABASE_USERNAME"),
			MUTPassword:        os.Getenv("MUT_DATABASE_PASSWORD"),
			MUTDBName:          os.Getenv("MUT_DATABASE_NAME"),
			MUTSchema:          os.Getenv("MUT_DATABASE_SCHEMA"),
			MUTMaxIdle:         cast.ToInt(os.Getenv("MUT_DATABASE_MAX_IDLE")),
			MUTMaxConn:         cast.ToInt(os.Getenv("MUT_DATABASE_MAX_CONN")),
			MUTConnMaxLifetime: cast.ToInt(os.Getenv("MUT_DATABASE_CONN_LIFETIME")),
		},
		Cache: Redis{
			Host:         os.Getenv("REDIS_HOST"),
			Port:         cast.ToInt(os.Getenv("REDIS_PORT")),
			Username:     os.Getenv("REDIS_USERNAME"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           cast.ToInt(os.Getenv("REDIS_DB")),
			UseTLS:       cast.ToBool(os.Getenv("REDIS_USE_TLS")),
			MaxRetries:   cast.ToInt(os.Getenv("REDIS_MAX_RETRIES")),
			MinIdleConns: cast.ToInt(os.Getenv("REDIS_MIN_IDLE_CONNS")),
			PoolSize:     cast.ToInt(os.Getenv("REDIS_POOL_SIZE")),
			PoolTimeout:  cast.ToInt(os.Getenv("REDIS_POOL_TIMEOUT")),
			MaxConnAge:   cast.ToInt(os.Getenv("REDIS_MAX_CONN_AGE")),
			ReadTimeout:  cast.ToInt(os.Getenv("REDIS_READ_TIMEOUT")),
			WriteTimeout: cast.ToInt(os.Getenv("REDIS_WRITE_TIMEOUT")),
		},
	}

	return
}
