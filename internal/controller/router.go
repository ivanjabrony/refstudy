package controller

import (
	"ivanjabrony/refstudy/docs"
	"ivanjabrony/refstudy/internal/logger"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(logger *logger.MyLogger, userUsecase UserUsecase, validator *validator.Validate) *gin.Engine {
	r := gin.Default()

	// timeoutTime := os.Getenv("TIMEOUT_TIME")
	// if timeoutTime == "" {
	// 	timeoutTime = "3"
	// }
	// timeoutParsed, err := strconv.Atoi(timeoutTime)
	// if err != nil {
	// 	timeoutParsed = 3
	// }

	// r.Use(middleware.LoggerMiddleware(logger))
	// r.Use(middleware.TimeoutMiddleware(time.Duration(timeoutParsed) * time.Second))

	userCotroller := NewUserController(userUsecase, validator)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/users")

	api.POST("/", userCotroller.CreateUser)
	api.PUT("/", userCotroller.UpdateUser)
	api.GET("/:id", userCotroller.GetUser)
	api.DELETE("/:id", userCotroller.DeleteUserById)
	api.GET("/", userCotroller.GetAllUsers)

	return r
}
