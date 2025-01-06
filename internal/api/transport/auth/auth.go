package auth

import (
	"github.com/bigxxby/dream-test-task/internal/api/service/auth"
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthCtrl represents the Auth Controller
type AuthCtrl struct {
	AuthService auth.IAuthService
}

// IAuthCtrl defines the Auth Controller interface
type IAuthCtrl interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Whoami(ctx *gin.Context)
}

// NewAuthController creates a new instance of AuthCtrl
func NewAuthController(service auth.IAuthService) IAuthCtrl {
	return &AuthCtrl{AuthService: service}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse defines the structure for register response body
type RegisterResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	User    models.User `json:"user"`
}

// User represents the user structure

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Register request body"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (ac AuthCtrl) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Message: "Bad request",
			Success: false,
		})
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.JSON(400, ErrorResponse{
			Error:   "Username or password is empty",
			Message: "Bad request",
			Success: false,
		})
		return
	}

	user, status, err := ac.AuthService.Register(req.Username, req.Password)
	if err != nil {
		switch status {
		case 400:
			ctx.JSON(400, ErrorResponse{
				Error:   err.Error(),
				Message: "Bad request",
				Success: false,
			})
			return
		case 500:
			ctx.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Message: "Internal server error",
				Success: false,
			})
			return
		case 409:
			ctx.JSON(409, ErrorResponse{
				Error:   err.Error(),
				Message: "User already exists",
				Success: false,
			})
			return
		default:
			ctx.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Message: "Internal server error",
				Success: false,
			})
			return
		}
	}

	ctx.JSON(200, RegisterResponse{
		User:    *user,
		Message: "User created",
		Success: true,
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login request body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (ac AuthCtrl) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Message: "Bad request",
			Success: false,
		})
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.JSON(400, ErrorResponse{
			Error:   "Username or password is empty",
			Message: "Bad request",
			Success: false,
		})
		return
	}

	token, status, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		switch status {
		case 400:
			ctx.JSON(400, ErrorResponse{
				Error:   err.Error(),
				Message: "Bad request",
				Success: false,
			})
			return
		case 500:
			ctx.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Message: "Internal server error",
				Success: false,
			})
			return
		case 404:
			ctx.JSON(404, ErrorResponse{
				Error:   err.Error(),
				Message: "User not found",
				Success: false,
			})
			return
		case 401:
			ctx.JSON(401, ErrorResponse{
				Error:   err.Error(),
				Message: "Wrong password",
				Success: false,
			})
			return
		default:
			ctx.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Message: "Internal server error",
				Success: false,
			})
			return
		}
	}

	ctx.JSON(200, LoginResponse{
		Token:   token,
		Message: "User logged in",
		Success: true,
	})
}

// WhoamiResponse defines the structure for whoami response body
type WhoamiResponse struct {
	User    models.User `json:"user"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

// Whoami godoc
// @Summary      Get user information
// @Description  Returns information about the currently authenticated user
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} WhoamiResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/whoami [get]
func (ac AuthCtrl) Whoami(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if userID == nil {
		ctx.JSON(401, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Unauthorized",
			Success: false,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Message: "Bad request",
			Success: false,
		})
		return
	}

	user, status, err := ac.AuthService.WHOAMI(&userIDUUID)
	if err != nil {
		switch status {
		case 404:
			ctx.JSON(404, ErrorResponse{
				Error:   err.Error(),
				Message: "User not found",
				Success: false,
			})
			return
		default:
			ctx.JSON(500, ErrorResponse{
				Error:   err.Error(),
				Message: "Internal server error",
				Success: false,
			})
			return
		}
	}

	ctx.JSON(200, WhoamiResponse{
		User:    *user,
		Message: "User found",
		Success: true,
	})
}
