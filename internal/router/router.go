package router

import (
	"github.com/bigxxby/dream-test-task/internal/api/middleware"
	authRepo "github.com/bigxxby/dream-test-task/internal/api/repo/auth"
	authService "github.com/bigxxby/dream-test-task/internal/api/service/auth"
	authController "github.com/bigxxby/dream-test-task/internal/api/transport/auth"

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

	//create group
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/whoami", middleware.AuthMiddleware(), authController.Whoami)

	}

	return router, nil

}
