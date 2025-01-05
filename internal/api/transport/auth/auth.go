package auth

import (
	"github.com/bigxxby/dream-test-task/internal/api/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthCtrl struct {
	AuthService auth.IAuthService
}

type IAuthCtrl interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Whoami(ctx *gin.Context)
}

func NewAuthController(service auth.IAuthService) IAuthCtrl {
	return &AuthCtrl{AuthService: service}
}

func (ac AuthCtrl) Register(ctx *gin.Context) {
	type registerRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req registerRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Bad request",
			"success": false,
		})
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.JSON(400, gin.H{
			"error":   "Username or password is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	user, status, err := ac.AuthService.Register(req.Username, req.Password)
	if err != nil {
		switch status {
		case 400:
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "Bad request",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		case 409:
			ctx.JSON(409, gin.H{
				"error":   err.Error(),
				"message": "User already exists",
				"success": false,
			})
			return
		default:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}

	}

	ctx.JSON(200, gin.H{
		"user":    user,
		"message": "User created",
		"success": true,
	})

}

func (ac AuthCtrl) Login(ctx *gin.Context) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.JSON(400, gin.H{
			"error":   "Username or password is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	token, status, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		switch status {
		case 400:
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "Bad request",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		case 404:
			ctx.JSON(404, gin.H{
				"error":   err.Error(),
				"message": "User not found",
				"success": false,
			})
			return
		case 401:
			ctx.JSON(401, gin.H{
				"error":   err.Error(),
				"message": "Wrong password",
				"success": false,
			})
			return
		default:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"token":   token,
		"message": "User logged in",
		"success": true,
	})
}

func (ac AuthCtrl) Whoami(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if userID == nil {
		ctx.JSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Bad request",
			"success": false,
		})
		return
	}

	user, status, err := ac.AuthService.WHOAMI(&userIDUUID)
	if err != nil {
		switch status {
		case 404:
			ctx.JSON(404, gin.H{
				"error":   err.Error(),
				"message": "User not found",
				"success": false,
			})
			return
		default:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"user":    user,
		"message": "User found",
		"success": true,
	})

}
