package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nabd/controllers"
	"nabd/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthController_NewAuthController(t *testing.T) {
	config := &models.Config{}
	config.Auth.AdminToken = "test-token"
	
	controller := controllers.NewAuthController(config)
	
	assert.NotNil(t, controller)
}

func TestAuthController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &models.Config{}
	config.Auth.AdminToken = "valid-admin-token"
	
	controller := controllers.NewAuthController(config)
	
	requestBody := map[string]string{
		"token": "valid-admin-token",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	
	router := gin.New()
	router.POST("/login", controller.Login)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "message")
	assert.Equal(t, "Authentication successful", response["message"])
	assert.NotEmpty(t, response["token"])
}

func TestAuthController_Login_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &models.Config{}
	config.Auth.AdminToken = "valid-admin-token"
	
	controller := controllers.NewAuthController(config)
	
	requestBody := map[string]string{
		"token": "invalid-token",
	}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	
	router := gin.New()
	router.POST("/login", controller.Login)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid admin token")
}

func TestAuthController_Login_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &models.Config{}
	config.Auth.AdminToken = "valid-admin-token"
	
	controller := controllers.NewAuthController(config)
	
	requestBody := map[string]string{}
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	
	router := gin.New()
	router.POST("/login", controller.Login)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Login_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &models.Config{}
	config.Auth.AdminToken = "valid-admin-token"
	
	controller := controllers.NewAuthController(config)
	
	router := gin.New()
	router.POST("/login", controller.Login)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &models.Config{}
	config.Auth.AdminToken = "valid-admin-token"
	
	controller := controllers.NewAuthController(config)
	
	router := gin.New()
	router.POST("/login", controller.Login)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}