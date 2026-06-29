package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTaskTestRouter(t *testing.T) (*gin.Engine, uint) {
	// Setup DB
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.DB = db
	database.DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Task{},
		&models.Tag{},
		&models.Comment{},
	)

	// Create test user
	user := models.User{
		Name:     "Task User",
		Email:    "task@test.com",
		Password: "hashedpass",
	}
	database.DB.Create(&user)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Auth middleware mock
	r.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Set("user", user)
		c.Next()
	})

	r.GET("/api/tasks", ListTasks)
	r.POST("/api/tasks", CreateTask)
	r.GET("/api/tasks/:id", GetTask)
	r.PUT("/api/tasks/:id", UpdateTask)
	r.DELETE("/api/tasks/:id", DeleteTask)

	return r, user.ID
}

func TestCreateTask(t *testing.T) {
	router, userID := setupTaskTestRouter(t)

	// Create workspace
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	payload := TaskReq{
		Title:       "Buy groceries",
		Description: "Milk, eggs, bread",
		Status:      "todo",
		Priority:    "high",
		WorkspaceID: ws.ID,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["data"])
}

func TestListTasks(t *testing.T) {
	router, userID := setupTaskTestRouter(t)

	// Create workspace
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	// Create task
	task := models.Task{
		Title:       "Test task",
		Description: "Test desc",
		Status:      "todo",
		Priority:    "medium",
		WorkspaceID: ws.ID,
		UserID:      userID,
	}
	database.DB.Create(&task)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/tasks?workspace_id=%d", ws.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["data"])
	assert.Equal(t, float64(1), response["total"])
}

func TestGetTask(t *testing.T) {
	router, userID := setupTaskTestRouter(t)

	// Create workspace + task
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	task := models.Task{
		Title:       "Get me",
		Description: "Test",
		Status:      "todo",
		Priority:    "low",
		WorkspaceID: ws.ID,
		UserID:      userID,
	}
	database.DB.Create(&task)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Get me", data["title"])
}

func TestUpdateTask(t *testing.T) {
	router, userID := setupTaskTestRouter(t)

	// Create workspace + task
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	task := models.Task{
		Title:       "Old title",
		Description: "Test",
		Status:      "todo",
		Priority:    "medium",
		WorkspaceID: ws.ID,
		UserID:      userID,
	}
	database.DB.Create(&task)

	payload := TaskReq{
		Title:    "Updated title",
		Status:   "doing",
		Priority: "high",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/tasks/%d", task.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Updated title", data["title"])
	assert.Equal(t, "doing", data["status"])
}

func TestDeleteTask(t *testing.T) {
	router, userID := setupTaskTestRouter(t)

	// Create workspace + task
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	task := models.Task{
		Title:       "Delete me",
		Status:      "todo",
		Priority:    "low",
		WorkspaceID: ws.ID,
		UserID:      userID,
	}
	database.DB.Create(&task)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify deleted
	var count int64
	database.DB.Model(&models.Task{}).Where("id = ?", task.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
