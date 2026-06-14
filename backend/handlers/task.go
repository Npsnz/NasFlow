package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskReq represents the payload to create/update tasks
type TaskReq struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	WorkspaceID  uint       `json:"workspace_id"`
	DueDate      *time.Time `json:"due_date"`
	IsRecurring  *bool      `json:"is_recurring"`
	RecurRule    string     `json:"recur_rule"`
	ParentTaskID *uint      `json:"parent_task_id"`
	SortOrder    *float64   `json:"sort_order"`
	TagIDs       []uint     `json:"tag_ids"`
}

func ListTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Filter parameters
	workspaceIDStr := c.Query("workspace_id")
	status := c.Query("status")
	priority := c.Query("priority")
	tagIDStr := c.Query("tag_id")
	dueBeforeStr := c.Query("due_before")
	dueAfterStr := c.Query("due_after")
	search := c.Query("search")
	parentOnly := c.Query("parent_only") == "true"

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	query := database.DB.Preload("Tags").Preload("Subtasks").Where("user_id = ?", userID)

	if workspaceIDStr != "" {
		if wID, err := strconv.ParseUint(workspaceIDStr, 10, 32); err == nil {
			query = query.Where("workspace_id = ?", wID)
		}
	}

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// Don't show archived tasks by default in general lists
		query = query.Where("status != ?", "archived")
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if parentOnly {
		query = query.Where("parent_task_id IS NULL")
	}

	if dueBeforeStr != "" {
		if t, err := time.Parse(time.RFC3339, dueBeforeStr); err == nil {
			query = query.Where("due_date <= ?", t)
		}
	}

	if dueAfterStr != "" {
		if t, err := time.Parse(time.RFC3339, dueAfterStr); err == nil {
			query = query.Where("due_date >= ?", t)
		}
	}

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if tagIDStr != "" {
		if tID, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			query = query.Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
				Where("task_tags.tag_id = ?", tID)
		}
	}

	// Count total matching
	var total int64
	query.Model(&models.Task{}).Count(&total)

	// Sort order: LexoRank SortOrder asc, then CreatedAt desc
	query = query.Order("sort_order ASC, created_at DESC")

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถดึงข้อมูลงานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": total,
	})
}

func CreateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "กรุณาระบุหัวข้องาน",
			"code":  "TITLE_REQUIRED",
		})
		return
	}

	// Calculate SortOrder if not set
	var sortOrder float64
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	} else {
		// Default to Max(sort_order) + 1000.0 within the workspace & status column
		var maxSort float64
		database.DB.Model(&models.Task{}).
			Where("user_id = ? AND workspace_id = ? AND status = ?", userID, req.WorkspaceID, req.Status).
			Select("COALESCE(MAX(sort_order), 0)").Row().Scan(&maxSort)
		sortOrder = maxSort + 1000.0
	}

	isRecur := false
	if req.IsRecurring != nil {
		isRecur = *req.IsRecurring
	}

	task := models.Task{
		Title:        req.Title,
		Description:  req.Description,
		Status:       req.Status,
		Priority:     req.Priority,
		WorkspaceID:  req.WorkspaceID,
		UserID:       userID,
		DueDate:      req.DueDate,
		IsRecurring:  isRecur,
		RecurRule:    req.RecurRule,
		ParentTaskID: req.ParentTaskID,
		SortOrder:    sortOrder,
	}

	if task.Status == "" {
		task.Status = "todo"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	tx := database.DB.Begin()
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถสร้างงานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Associate tags if passed
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ? AND user_id = ?", req.TagIDs, userID).Find(&tags).Error; err == nil && len(tags) > 0 {
			tx.Model(&task).Association("Tags").Replace(tags)
		}
	}
	tx.Commit()

	// Preload relations to return full object
	database.DB.Preload("Tags").Preload("Subtasks").First(&task, task.ID)

	// Broadcast SSE Event
	BroadcastEvent(userID, "task.created", task)

	c.JSON(http.StatusCreated, gin.H{
		"data": task,
	})
}

func GetTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var task models.Task
	if err := database.DB.Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Preload("Subtasks").
		Where("id = ? AND user_id = ?", id, userID).
		First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบงานที่กำหนด",
			"code":  "TASK_NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func UpdateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบงานที่ต้องการแก้ไข",
			"code":  "TASK_NOT_FOUND",
		})
		return
	}

	var req TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	statusChanged := false
	oldStatus := task.Status

	if req.Title != "" {
		task.Title = req.Title
	}
	task.Description = req.Description

	if req.Status != "" {
		if task.Status != req.Status {
			statusChanged = true
			task.Status = req.Status
			if req.Status == "done" {
				now := time.Now()
				task.CompletedAt = &now
			} else {
				task.CompletedAt = nil
			}
		}
	}

	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.WorkspaceID != 0 {
		task.WorkspaceID = req.WorkspaceID
	}
	task.DueDate = req.DueDate
	task.ParentTaskID = req.ParentTaskID

	if req.IsRecurring != nil {
		task.IsRecurring = *req.IsRecurring
	}
	task.RecurRule = req.RecurRule

	if req.SortOrder != nil {
		task.SortOrder = *req.SortOrder
	}

	tx := database.DB.Begin()
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกการแก้ไขงานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Update tags association if specified
	if req.TagIDs != nil {
		var tags []models.Tag
		if len(req.TagIDs) > 0 {
			if err := tx.Where("id IN ? AND user_id = ?", req.TagIDs, userID).Find(&tags).Error; err == nil {
				tx.Model(&task).Association("Tags").Replace(tags)
			}
		} else {
			tx.Model(&task).Association("Tags").Clear()
		}
	}
	tx.Commit()

	// Load full updated object
	database.DB.Preload("Tags").Preload("Subtasks").First(&task, task.ID)

	// Broadcast SSE Event
	if statusChanged {
		BroadcastEvent(userID, "task.moved", map[string]interface{}{
			"task":       task,
			"old_status": oldStatus,
		})
	} else {
		BroadcastEvent(userID, "task.updated", task)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func DeleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบงานที่ต้องการลบ",
			"code":  "TASK_NOT_FOUND",
		})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถลบงานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Broadcast SSE Event
	BroadcastEvent(userID, "task.deleted", map[string]interface{}{
		"id": id,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "ลบงานเรียบร้อยแล้ว",
			"id":      id,
		},
	})
}

func UpdateTaskStatus(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	type StatusReq struct {
		Status string `json:"status" binding:"required"`
	}

	var req StatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบงาน",
			"code":  "TASK_NOT_FOUND",
		})
		return
	}

	oldStatus := task.Status
	task.Status = req.Status
	if req.Status == "done" {
		now := time.Now()
		task.CompletedAt = &now
	} else {
		task.CompletedAt = nil
	}

	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกสถานะได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Preload tags & subtasks
	database.DB.Preload("Tags").Preload("Subtasks").First(&task, task.ID)

	// Broadcast SSE Event
	BroadcastEvent(userID, "task.moved", map[string]interface{}{
		"task":       task,
		"old_status": oldStatus,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func ReorderTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	type ReorderItem struct {
		ID        uint    `json:"id"`
		SortOrder float64 `json:"sort_order"`
		Status    string  `json:"status"` // Optional status change on drag-and-drop
	}

	type ReorderReq struct {
		Tasks []ReorderItem `json:"tasks" binding:"required"`
	}

	var req ReorderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	tx := database.DB.Begin()
	for _, item := range req.Tasks {
		updateMap := map[string]interface{}{
			"sort_order": item.SortOrder,
		}
		if item.Status != "" {
			updateMap["status"] = item.Status
			if item.Status == "done" {
				updateMap["completed_at"] = time.Now()
			} else {
				updateMap["completed_at"] = nil
			}
		}

		if err := tx.Model(&models.Task{}).Where("id = ? AND user_id = ?", item.ID, userID).Updates(updateMap).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "การจัดลำดับล้มเหลว",
				"code":  "DATABASE_ERROR",
			})
			return
		}
	}
	tx.Commit()

	// Broadcast updated event for general tasks list refresh
	BroadcastEvent(userID, "task.updated", gin.H{"reordered": true})

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "จัดเรียงงานสำเร็จแล้ว",
		},
	})
}

func CompleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบงาน",
			"code":  "TASK_NOT_FOUND",
		})
		return
	}

	now := time.Now()
	task.Status = "done"
	task.CompletedAt = &now

	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถทำเครื่องหมายว่าเสร็จงานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Trigger recurring task creation if applicable
	var nextTask *models.Task
	if task.IsRecurring {
		var err error
		nextTask, err = services.CreateNextOccurrence(database.DB, task)
		if err != nil {
			// Log error, but don't fail completion
			println("Error creating recurring task occurrence: ", err.Error())
		}
	}

	// Preload tags & subtasks
	database.DB.Preload("Tags").Preload("Subtasks").First(&task, task.ID)

	// Broadcast SSE Event for current task completed
	BroadcastEvent(userID, "task.moved", map[string]interface{}{
		"task":       task,
		"old_status": "todo",
	})

	response := gin.H{
		"data": task,
	}

	// If next recurring occurrence is created, preload and broadcast it
	if nextTask != nil {
		database.DB.Preload("Tags").Preload("Subtasks").First(nextTask, nextTask.ID)
		BroadcastEvent(userID, "task.created", nextTask)
		response["next_occurrence"] = nextTask
	}

	c.JSON(http.StatusOK, response)
}

func GetOverdue(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	tasks, err := services.GetOverdueTasks(database.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถโหลดงานที่เกินกำหนดได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": len(tasks),
	})
}

func GetStats(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var total, todo, doing, doneToday, overdue, dueThisWeek int64
	now := time.Now()

	// Total active tasks (not archived)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status != 'archived'", userID).Count(&total)

	// Todo status
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = 'todo'", userID).Count(&todo)

	// Doing status
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = 'doing'", userID).Count(&doing)

	// Done completed today
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = 'done' AND completed_at >= ?", userID, todayStart).Count(&doneToday)

	// Overdue
	database.DB.Model(&models.Task{}).Where("user_id = ? AND due_date < ? AND status != 'done' AND status != 'archived'", userID, now).Count(&overdue)

	// Due this week (between start of week e.g. Monday and end of week e.g. Sunday)
	weekdayOffset := int(now.Weekday())
	if weekdayOffset == 0 {
		weekdayOffset = 7 // Adjust Sunday to be day 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -weekdayOffset+1)
	weekEnd := weekStart.AddDate(0, 0, 7)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND due_date >= ? AND due_date < ? AND status != 'done' AND status != 'archived'", userID, weekStart, weekEnd).Count(&dueThisWeek)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total":         total,
			"todo":          todo,
			"doing":         doing,
			"done_today":    doneToday,
			"overdue":       overdue,
			"due_this_week": dueThisWeek,
		},
	})
}

// AddComment inserts a comment for a task
func AddComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสงานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	type CommentReq struct {
		Content string `json:"content" binding:"required"`
	}

	var req CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "กรุณาระบุข้อความความคิดเห็น",
			"code":  "CONTENT_REQUIRED",
		})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		TaskID:  uint(taskID),
		UserID:  userID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกความคิดเห็นได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Preload the User details
	database.DB.Preload("User").First(&comment, comment.ID)

	// Fetch full task to broadcast updated comments state
	var task models.Task
	if err := database.DB.Preload("Tags").Preload("Subtasks").First(&task, taskID).Error; err == nil {
		BroadcastEvent(userID, "task.updated", task)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": comment,
	})
}

// DeleteComment deletes comment owned by user
func DeleteComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสความคิดเห็นไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบความคิดเห็นที่ต้องการลบหรือไม่มีสิทธิ์",
			"code":  "COMMENT_NOT_FOUND",
		})
		return
	}

	taskID := comment.TaskID
	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถลบความคิดเห็นได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Fetch full task to broadcast updated comments state
	var task models.Task
	if err := database.DB.Preload("Tags").Preload("Subtasks").First(&task, taskID).Error; err == nil {
		BroadcastEvent(userID, "task.updated", task)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "ลบความคิดเห็นสำเร็จแล้ว",
			"id":      commentID,
		},
	})
}
