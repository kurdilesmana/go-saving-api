package deps

import (
	"log"

	config "github.com/kurdilesmana/go-saving-api/apps/account/configs"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/adapters/v1/repositories/accountRepo"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/adapters/v1/repositories/cacheRepo"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/services/accountService"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/services/transactionService"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/infra/db"
	middle "github.com/kurdilesmana/go-saving-api/apps/account/server/middlewares"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/sirupsen/logrus"
)

const (
	keyTransaction = "account-ctx"
	timeout        = 60
)

type Dependency struct {
	Cfg                config.EnvironmentConfig
	AccountService     accountPort.IAccountService
	TransactionService transactionPort.ITransactionService
	AuthMiddleware     *middle.AuthMiddleware
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
	database, err := db.OpenPgsqlConnection(&config.ACCDatabase, logger)
	if err != nil {
		log.Panic(err)
	}

	redis, err := db.RedisNewClient(&config.Cache, logger)
	if err != nil {
		log.Panic(err)
	}

	// BIG DEPENDENCY STAGE END =======================================

	// init repository
	accountRepository := accountRepo.NewaccountRepo(database, keyTransaction, timeout, logger)
	cacheRepository := cacheRepo.NewCacheRepo(redis, keyTransaction, timeout, logger)

	//init middleware
	authMiddleware := middle.NewAuthMiddleware(accountRepository, logger)

	// init service
	accountService := accountService.NewAccountService(accountRepository, logger)
	transactionService := transactionService.NewTransactionService(accountRepository, cacheRepository, logger)

	return Dependency{
		Cfg:                config,
		AccountService:     accountService,
		TransactionService: transactionService,
		AuthMiddleware:     authMiddleware,
		Validator:          validator,
		Logger:             logger,
	}
}
