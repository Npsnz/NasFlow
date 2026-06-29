package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	// Initialize test database (SQLite in-memory)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = db

	// Migrate models
	err = database.DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Task{},
		&models.Tag{},
		&models.Comment{},
	)
	assert.NoError(t, err)
}

func setupTestRouter(t *testing.T) *gin.Engine {
	setupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/auth/register", Register)
	r.POST("/api/auth/login", Login)
	return r
}

func TestRegisterUser(t *testing.T) {
	router := setupTestRouter(t)

	payload := RegisterReq{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["token"])
}

func TestLoginUser(t *testing.T) {
	router := setupTestRouter(t)

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create user directly
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	database.DB.Create(&user)

	payload := LoginReq{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["token"])
}

func TestRegisterDuplicateEmail(t *testing.T) {
	router := setupTestRouter(t)

	payload := RegisterReq{
		Name:     "User 1",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Try same email again
	payload2 := RegisterReq{
		Name:     "User 2",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}
