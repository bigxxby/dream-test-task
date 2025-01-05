package router

import (
	"github.com/bigxxby/dream-test-task/internal/api/middleware"
	authRepo "github.com/bigxxby/dream-test-task/internal/api/repo/auth"
	authService "github.com/bigxxby/dream-test-task/internal/api/service/auth"
	authController "github.com/bigxxby/dream-test-task/internal/api/transport/auth"

	shortenerRepo "github.com/bigxxby/dream-test-task/internal/api/repo/shortener"
	shortenerService "github.com/bigxxby/dream-test-task/internal/api/service/shortener"
	shortenerController "github.com/bigxxby/dream-test-task/internal/api/transport/shortener"

	userRepo "github.com/bigxxby/dream-test-task/internal/api/repo/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) (*gin.Engine, error) {
	router := gin.Default()

	userRepo := userRepo.NewUserRepo(db)
	//

	authRepo := authRepo.NewAuthRepo(db)
	authService := authService.NewAuthService(authRepo, userRepo)
	authController := authController.NewAuthController(authService)

	shortenerRepo := shortenerRepo.NewShortenerRepo(db)
	shortenerService := shortenerService.NewShortenerService(shortenerRepo)
	shortenerController := shortenerController.NewShortenerController(shortenerService)

	//create group
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/whoami", middleware.AuthMiddleware(), authController.Whoami)

	}
	shortener := router.Group("/shortener")
	{
		shortener.GET("/", middleware.AuthMiddleware(), shortenerController.GetLinks)
		shortener.GET("/:shortID", middleware.AuthMiddleware(), shortenerController.Redirect)
		shortener.GET("/stats/:shortID", middleware.AuthMiddleware(), shortenerController.GetLink)
		shortener.POST("/", middleware.AuthMiddleware(), shortenerController.CreateShortLink)
		shortener.DELETE("/:shortID", middleware.AuthMiddleware(), shortenerController.DeleteLink)

	}

	return router, nil

}
