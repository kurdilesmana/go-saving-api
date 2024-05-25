package deps

import (
	"log"

	config "github.com/kurdilesmana/go-saving-api/apps/mutation/configs"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/adapters/v1/repositories/cacheRepo"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/adapters/v1/repositories/transactionRepo"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/services/transactionService"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/infra/db"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	keyTransaction = "mutation-ctx"
	timeout        = 60
)

type Dependency struct {
	Cfg                config.EnvironmentConfig
	TransactionService transactionPort.ITransactionService
	Redis              *redis.Client
	Validator          *validator.RequestValidator
	Logger             *logging.Logger
}

func SetupDependencies() Dependency {
	// init validator
	validator := validator.NewRequestValidator()

	// init config
	config, err := config.LoadENVConfig()
	if err != nil {
		log.Panic(err)
	}

	// load logger
	logger := logging.NewLogger(config.AppConfig.Name)
	defer logger.Info(logrus.Fields{}, nil, "Done cleanup tasks...")

	// BIG DEPENDENCY STAGE =======================================
	database, err := db.OpenPgsqlConnection(&config.MUTDatabase, logger)
	if err != nil {
		log.Panic(err)
	}

	redis, err := db.RedisNewClient(&config.Cache, logger)
	if err != nil {
		log.Panic(err)
	}

	// BIG DEPENDENCY STAGE END =======================================

	// init repository
	transactionRepository := transactionRepo.NewMutationRepo(database, keyTransaction, timeout, logger)
	cacheRepository := cacheRepo.NewCacheRepo(redis, keyTransaction, timeout, logger)

	// init service
	transactionService := transactionService.NewTransactionService(transactionRepository, cacheRepository, logger)

	return Dependency{
		Cfg:                config,
		TransactionService: transactionService,
		Redis:              redis,
		Validator:          validator,
		Logger:             logger,
	}
}
