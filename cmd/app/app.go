package app

import (
	"ivanjabrony/refstudy/internal/controller"
	"ivanjabrony/refstudy/internal/logger"
	"ivanjabrony/refstudy/internal/repository"
	"ivanjabrony/refstudy/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router *gin.Engine
	db     *pgxpool.Pool
}

func New(db *pgxpool.Pool) *App {
	logger := logger.New(getLogLevel(), logger.LogFormatText)
	repositories := mustInitRepositories(db, logger)
	usecases := mustInitUsecases(repositories, logger)
	validator := validator.New(validator.WithRequiredStructEnabled())

	router := controller.SetupRouter(
		logger,
		usecases.user,
		validator,
	)

	return &App{
		Router: router,
		db:     db,
	}
}

func (a *App) Run(addr string) error {
	return a.Router.Run(addr)
}

type repositories struct {
	user *repository.UserRepository
}

type usecases struct {
	user *usecase.UserUsecase
}

func mustInitRepositories(db *pgxpool.Pool, logger *logger.MyLogger) *repositories {
	repository, err := repository.NewUserRepository(db, logger)
	if err != nil {
		panic(err)
	}
	return &repositories{
		user: repository,
	}
}

func mustInitUsecases(r *repositories, logger *logger.MyLogger) *usecases {
	if r == nil || logger == nil {
		log.Fatal("couldn't init usecases: nil values in constructor")
	}
	user, err := usecase.NewUserUsecase(r.user, logger)
	if err != nil {
		log.Fatalf("couldn't init usecases: %w", err)
	}

	return &usecases{user: user}
}

func getLogLevel() logger.LoggerLevel {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return logger.Debug
	case "prod":
		return logger.Prod
	case "test":
		return logger.Test
	default:
		return logger.Debug
	}
}
