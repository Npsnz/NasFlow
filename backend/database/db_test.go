package database

import (
	"os"
	"testing"

	"backend/config"
	"backend/models"
)

func TestInitDBAndSeeding(t *testing.T) {
	// 1. Setup config for testing
	config.AppConfig.DBPath = "test_taskflow.db"
	defer os.Remove(config.AppConfig.DBPath) // Cleanup DB file after test

	// 2. Initialize Database and run migrations
	InitDB()
	if DB == nil {
		t.Fatal("Expected DB connection, got nil")
	}

	// 3. Create dummy user
	user := models.User{
		Name:     "Test User",
		Email:    "test@taskflow.local",
		Password: "hashedpassword123",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create dummy user: %v", err)
	}

	// 4. Seed default workspaces
	err := SeedDefaultWorkspaces(user.ID)
	if err != nil {
		t.Fatalf("Failed to seed default workspaces: %v", err)
	}

	// 5. Verify workspaces were created
	var workspaces []models.Workspace
	if err := DB.Where("user_id = ?", user.ID).Order("sort_order ASC").Find(&workspaces).Error; err != nil {
		t.Fatalf("Failed to query seeded workspaces: %v", err)
	}

	if len(workspaces) != 3 {
		t.Errorf("Expected 3 workspaces, got %d", len(workspaces))
	}

	expectedNames := []string{"งาน", "ส่วนตัว", "สุขภาพ"}
	expectedIcons := []string{"briefcase", "home", "heart"}
	expectedColors := []string{"#534AB7", "#1D9E75", "#D85A30"}

	for i, ws := range workspaces {
		if ws.Name != expectedNames[i] {
			t.Errorf("Index %d: expected name %q, got %q", i, expectedNames[i], ws.Name)
		}
		if ws.Icon != expectedIcons[i] {
			t.Errorf("Index %d: expected icon %q, got %q", i, expectedIcons[i], ws.Icon)
		}
		if ws.Color != expectedColors[i] {
			t.Errorf("Index %d: expected color %q, got %q", i, expectedColors[i], ws.Color)
		}
	}
}
