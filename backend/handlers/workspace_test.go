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

func setupWorkspaceTestRouter(t *testing.T) (*gin.Engine, uint) {
	// Setup DB
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.DB = db
	database.DB.AutoMigrate(&models.User{}, &models.Workspace{}, &models.Task{}, &models.Tag{}, &models.Comment{})

	// Create test user
	user := models.User{
		Name:     "Test User",
		Email:    "workspace@test.com",
		Password: "hashedpass",
	}
	database.DB.Create(&user)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add auth middleware mock
	r.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Set("user", user)
		c.Next()
	})

	r.GET("/api/workspaces", ListWorkspaces)
	r.POST("/api/workspaces", CreateWorkspace)
	r.PUT("/api/workspaces/:id", UpdateWorkspace)
	r.DELETE("/api/workspaces/:id", DeleteWorkspace)

	return r, user.ID
}

func TestListWorkspaces(t *testing.T) {
	router, userID := setupWorkspaceTestRouter(t)

	// Create workspace
	ws := models.Workspace{
		Name:   "Work",
		Slug:   "work-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	req := httptest.NewRequest("GET", "/api/workspaces", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["data"])
	assert.Equal(t, float64(1), response["total"])
}

func TestCreateWorkspace(t *testing.T) {
	router, _ := setupWorkspaceTestRouter(t)

	payload := WorkspaceReq{
		Name:  "Personal",
		Color: "#A89968",
		Icon:  "home",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/workspaces", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.NotNil(t, response["data"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Personal", data["name"])
}

func TestUpdateWorkspace(t *testing.T) {
	router, userID := setupWorkspaceTestRouter(t)

	// Create workspace
	ws := models.Workspace{
		Name:   "Old Name",
		Slug:   "old-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	payload := WorkspaceReq{
		Name:  "New Name",
		Color: "#8B7355",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/workspaces/%d", ws.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "New Name", data["name"])
	assert.Equal(t, "#8B7355", data["color"])
}

func TestDeleteWorkspace(t *testing.T) {
	router, userID := setupWorkspaceTestRouter(t)

	// Create workspace
	ws := models.Workspace{
		Name:   "Delete Me",
		Slug:   "delete-1",
		Color:  "#A89968",
		UserID: userID,
	}
	database.DB.Create(&ws)

	// Soft delete (default behavior)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/workspaces/%d", ws.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify soft-deleted (IsArchived = true)
	var softDeletedWs models.Workspace
	database.DB.Where("id = ?", ws.ID).First(&softDeletedWs)
	assert.True(t, softDeletedWs.IsArchived)
}
