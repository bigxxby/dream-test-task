package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bigxxby/dream-test-task/internal/api/transport/auth"
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the IAuthService interface.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, password string) (*models.User, int, error) {
	args := m.Called(username, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *MockAuthService) Login(username, password string) (string, int, error) {
	args := m.Called(username, password)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockAuthService) WHOAMI(userID *uuid.UUID) (*models.User, int, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func TestRegister(t *testing.T) {
	mockAuthService := new(MockAuthService)
	router := gin.Default()
	authCtrl := &auth.AuthCtrl{AuthService: mockAuthService}
	router.POST("/register", authCtrl.Register)

	// Тест с валидными данными
	mockAuthService.On("Register", "testuser", "password").Return(nil, 200, nil)

	reqBody := `{"username":"testuser","password":"password"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User created")

	// Тест с пустыми данными
	reqBody = `{"username":"","password":""}`
	req, _ = http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Username or password is empty")
}
func TestLogin(t *testing.T) {
	mockAuthService := new(MockAuthService)
	router := gin.Default()
	authCtrl := &auth.AuthCtrl{AuthService: mockAuthService}
	router.POST("/login", authCtrl.Login)

	// Тест с валидными данными
	mockAuthService.On("Login", "testuser", "password").Return("token", 200, nil)

	reqBody := `{"username":"testuser","password":"password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	// Тест с пустыми данными
	reqBody = `{"username":"","password":""}`
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Username or password is empty")
}
func TestWhoami(t *testing.T) {
	mockAuthService := new(MockAuthService)
	router := gin.Default()
	authCtrl := &auth.AuthCtrl{AuthService: mockAuthService}
	router.GET("/whoami", authCtrl.Whoami)

	// Тест с валидным пользователем
	mockAuthService.On("WHOAMI", mock.AnythingOfType("*uuid.UUID")).Return(nil, 200, nil)

	req, _ := http.NewRequest("GET", "/whoami", nil)
	req.Header.Set("user_id", "1234-5678-91011-1213") // Имитация заголовка с user_id
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")

	// Тест с отсутствующим user_id в заголовке
	req, _ = http.NewRequest("GET", "/whoami", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}
