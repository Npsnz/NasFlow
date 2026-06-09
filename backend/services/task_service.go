package services

import (
	"fmt"
	"log"
	"time"

	"backend/models"

	"gorm.io/gorm"
)

// BroadcastCallback registers a callback to push SSE alerts from this service
var BroadcastCallback func(userID uint, event string, payload interface{})

// CreateNextOccurrence duplicates a recurring task with a newly calculated due date
func CreateNextOccurrence(db *gorm.DB, task models.Task) (*models.Task, error) {
	if !task.IsRecurring || task.RecurRule == "" {
		return nil, nil
	}

	var baseDate time.Time
	if task.DueDate != nil {
		baseDate = *task.DueDate
	} else {
		baseDate = time.Now()
	}

	var nextDue time.Time
	switch task.RecurRule {
	case "daily":
		nextDue = baseDate.AddDate(0, 0, 1)
	case "weekly":
		nextDue = baseDate.AddDate(0, 0, 7)
	case "monthly":
		nextDue = baseDate.AddDate(0, 1, 0)
	case "weekdays":
		nextDue = baseDate.AddDate(0, 0, 1)
		for nextDue.Weekday() == time.Saturday || nextDue.Weekday() == time.Sunday {
			nextDue = nextDue.AddDate(0, 0, 1)
		}
	default:
		return nil, fmt.Errorf("invalid recur rule: %s", task.RecurRule)
	}

	// Clone task properties to next task
	nextTask := models.Task{
		Title:        task.Title,
		Description:  task.Description,
		Status:       "todo",
		Priority:     task.Priority,
		WorkspaceID:  task.WorkspaceID,
		UserID:       task.UserID,
		DueDate:      &nextDue,
		IsRecurring:  true,
		RecurRule:    task.RecurRule,
		SortOrder:    task.SortOrder, // Use same LexoRank order as base
	}

	if err := db.Create(&nextTask).Error; err != nil {
		return nil, err
	}

	// Fetch tags of the original task to link them to the new occurrence
	var tags []models.Tag
	if err := db.Model(&task).Association("Tags").Find(&tags); err == nil && len(tags) > 0 {
		db.Model(&nextTask).Association("Tags").Replace(tags)
	}

	return &nextTask, nil
}

// GetOverdueTasks lists all tasks due in the past that are not completed or archived
func GetOverdueTasks(db *gorm.DB, userID uint) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	err := db.Preload("Tags").
		Where("user_id = ? AND due_date < ? AND status != 'done' AND status != 'archived'", userID, now).
		Find(&tasks).Error
	return tasks, err
}

// StartOverdueCheckService fires up an hourly routine to check and notify overdue tasks
func StartOverdueCheckService(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		// Run initial check on server start
		runOverdueCheck(db)
		for range ticker.C {
			runOverdueCheck(db)
		}
	}()
}

func runOverdueCheck(db *gorm.DB) {
	var userIDs []uint
	if err := db.Model(&models.User{}).Pluck("id", &userIDs).Error; err != nil {
		log.Printf("[Overdue Service] Error querying users: %v", err)
		return
	}

	for _, uID := range userIDs {
		tasks, err := GetOverdueTasks(db, uID)
		if err != nil {
			log.Printf("[Overdue Service] Error checking overdue tasks for user %d: %v", uID, err)
			continue
		}

		if len(tasks) > 0 && BroadcastCallback != nil {
			log.Printf("[Overdue Service] Alerting User %d with %d overdue tasks", uID, len(tasks))
			BroadcastCallback(uID, "tasks.overdue", map[string]interface{}{
				"count": len(tasks),
			})
		}
	}
}
