package video

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewVideoHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	baseRouter := router.Group("/api")

	// Test
	handler := NewVideoHandler(baseRouter)

	// Assert
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
}

func TestVideoHandler_Upload(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	baseRouter := router.Group("/api")
	handler := NewVideoHandler(baseRouter)
	handler.SetupRoutes()

	// Create a multipart form buffer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("video", "test.mp4")
	part.Write([]byte("test video content"))
	writer.Close()

	// Create test server
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/video/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assert response
	assert.NoError(t, err)
	assert.Equal(t, "Video upload successful", response["message"])
	assert.Equal(t, "test.mp4", response["filename"])
}
