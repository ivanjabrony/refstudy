package app

import (
	"ivanjabrony/refstudy/internal/controller"
	"ivanjabrony/refstudy/internal/repository"
	"ivanjabrony/refstudy/internal/usecase"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Router *gin.Engine
	db     *sqlx.DB
}

func New(db *sqlx.DB) *App {
	repositories := initRepositories(db)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel()}))
	usecases := initUsecases(repositories)

	router := controller.SetupRouter(
		logger,
		usecases.user,
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
	user repository.UserRepository
}

type usecases struct {
	user usecase.UserUsecase
}

func initRepositories(db *sqlx.DB) *repositories {
	return &repositories{
		user: repository.NewUserRepository(db),
	}
}

func initUsecases(r *repositories) *usecases {
	return &usecases{
		user: usecase.NewUserUsecase(r.user),
	}
}

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
