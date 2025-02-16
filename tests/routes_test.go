package tests

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/routes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMerchHandler struct {
	mock.Mock
}

func (m *MockMerchHandler) GetAllMerch(c *fiber.Ctx) error {
	return nil
}

func (m *MockMerchHandler) GetMerchByName(c *fiber.Ctx) error {
	return nil
}

// MockUserHandler мок для UserHandler
type MockUserHandler struct {
	mock.Mock
}

func (m *MockUserHandler) GetInfo(c *fiber.Ctx) error {
	return nil
}

func (m *MockUserHandler) SendCoins(c *fiber.Ctx) error {
	return nil
}

func (m *MockUserHandler) BuyItem(c *fiber.Ctx) error {
	return nil
}

func (m *MockUserHandler) AddCoins(c *fiber.Ctx) error {
	return nil
}

type MockAuthHandler struct {
	mock.Mock
}

func (m *MockAuthHandler) Authenticate(c *fiber.Ctx) error {
	return nil
}

func TestSetupRoutes(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "test-secret",
	}

	mockMerchHandler := &MockMerchHandler{}
	mockUserHandler := &MockUserHandler{}
	mockAuthHandler := &MockAuthHandler{}

	app := fiber.New()

	routes.SetupRoutes(cfg, app, mockMerchHandler, mockUserHandler, mockAuthHandler)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /api/merch",
			method:         "GET",
			path:           "/api/merch",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "GET /api/merch/:name",
			method:         "GET",
			path:           "/api/merch/test-item",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "POST /api/auth",
			method:         "POST",
			path:           "/api/auth",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "GET /api/info (protected)",
			method:         "GET",
			path:           "/api/info",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "POST /api/send-coins (protected)",
			method:         "POST",
			path:           "/api/send-coins",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "GET /api/buy/:item (protected)",
			method:         "GET",
			path:           "/api/buy/test-item",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "POST /api/add-coins/:amount (protected)",
			method:         "POST",
			path:           "/api/add-coins/100",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "GET /swagger/*",
			method:         "GET",
			path:           "/swagger/index.html",
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
