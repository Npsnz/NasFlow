package database

import (
	"fmt"
	"log"
	"backend/config"
	"backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Enable foreign keys for SQLite
	DB.Exec("PRAGMA foreign_keys = ON;")

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Task{},
		&models.Tag{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed successfully.")
}

// SeedDefaultWorkspaces seeds the 3 default workspaces for a new user:
// "งาน" (briefcase, #171717), "ส่วนตัว" (home, #1D9E75), "สุขภาพ" (heart, #D85A30)
func SeedDefaultWorkspaces(userID uint) error {
	defaults := []models.Workspace{
		{
			Name:      "งาน",
			Slug:      fmt.Sprintf("work-%d", userID),
			Color:     "#171717",
			Icon:      "briefcase",
			UserID:    userID,
			SortOrder: 1,
		},
		{
			Name:      "ส่วนตัว",
			Slug:      fmt.Sprintf("personal-%d", userID),
			Color:     "#1D9E75",
			Icon:      "home",
			UserID:    userID,
			SortOrder: 2,
		},
		{
			Name:      "สุขภาพ",
			Slug:      fmt.Sprintf("health-%d", userID),
			Color:     "#D85A30",
			Icon:      "heart",
			UserID:    userID,
			SortOrder: 3,
		},
	}

	for _, ws := range defaults {
		var count int64
		DB.Model(&models.Workspace{}).Where("user_id = ? AND name = ?", userID, ws.Name).Count(&count)
		if count == 0 {
			if err := DB.Create(&ws).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
