package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	baseRouter := router.Group("/api")

	// Test
	handler := NewUserHandler(baseRouter)

	// Assert
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
}

func TestUserHandler_SetupRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	baseRouter := router.Group("/api")
	handler := NewUserHandler(baseRouter)
	handler.SetupRoutes()

	// Create test server
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assert response
	assert.NoError(t, err)
	assert.Equal(t, "user routes", response["message"])
}
