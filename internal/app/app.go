package app

import (
	"fmt"

	"github.com/bulgil/blog-rest-api/internal/config"
	"github.com/bulgil/blog-rest-api/internal/logger"
	"github.com/bulgil/blog-rest-api/internal/router"
	"github.com/bulgil/blog-rest-api/internal/storage/psql"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Storage *pgxpool.Pool
	Router  *gin.Engine
}

func NewApp() *App {
	config := config.NewConfig()
	logger := logger.NewLogger(config.Env)
	storage := psql.NewStorage(config.PGStorage)
	router := router.NewRouter(config.Env, logger, storage)

	return &App{
		Config:  config,
		Logger:  logger,
		Storage: storage,
		Router:  router,
	}
}

func (app *App) Run() error {
	op := "app.Run"

	defer app.Storage.Close()
	app.Logger.Info("app started")
	if err := app.Router.Run(fmt.Sprintf("%s:%s",
		app.Config.HTTPServer.Address, app.Config.HTTPServer.Port)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
