package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"pvz/configs"
	"pvz/internal/logger"
	"pvz/internal/repositories"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

type Deps struct {
	AuthService      services.AuthService
	PvzService       services.PvzService
	ReceptionService services.ReceptionService
	ProductService   services.ProductService
}

func InitDeps() (Deps, error) {
	log := logger.Log.With("scope", "bootstrap")

	log.Info("initializing dependencies")

	dbDataSourceName, err := getDataSourceName()
	if err != nil {
		log.Error("failed to get data source name", "err", err)
		return Deps{}, err
	}

	db, err := sql.Open(configs.AppConfiguration.DB.DriverName, dbDataSourceName)
	if err != nil {
		log.Error("failed to open database connection", "err", err)
		return Deps{}, err
	}

	if err = db.Ping(); err != nil {
		log.Error("failed to ping database", "err", err)
		return Deps{}, err
	}
	log.Info("successfully connected to the database")

	log.Info("initializing repositories")
	userRepo := repositories.NewUserRepository(db)
	pvzRepo := repositories.NewPvzRepository(db)
	receptionRepo := repositories.NewReceptionRepository(db)
	productRepo := repositories.NewProductRepository(db)

	log.Info("initializing services")
	authService := services.NewAuthService(userRepo, db)
	pvzService := services.NewPvzService(pvzRepo, productRepo, receptionRepo, db)
	receptionService := services.NewReceptionService(receptionRepo, pvzRepo, db)
	productService := services.NewProductService(productRepo, pvzRepo, receptionRepo, db)

	log.Info("all dependencies initialized successfully")
	return Deps{
		AuthService:      authService,
		PvzService:       pvzService,
		ReceptionService: receptionService,
		ProductService:   productService,
	}, nil
}

func getDataSourceName() (string, error) {
	log := logger.Log.With("scope", "bootstrap", "func", "getDataSourceName")

	switch configs.AppConfiguration.DB.DriverName {
	case "postgres":
		log.Info("generating data source name for postgres")
		return fmt.Sprintf("postgres://%s:%s@%s:%v/%s",
			configs.AppConfiguration.DB.User,
			configs.AppConfiguration.DB.Password,
			configs.AppConfiguration.DB.Host,
			configs.AppConfiguration.DB.Port,
			configs.AppConfiguration.DB.Name,
		), nil
	default:
		log.Error("unsupported database driver", "driver", configs.AppConfiguration.DB.DriverName)
		return "", errors.NewInternalError()
	}
}
