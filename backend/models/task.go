package models

import (
	"time"
)

type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Status       string     `gorm:"default:'todo'" json:"status"` // todo | doing | done | archived
	Priority     string     `gorm:"default:'medium'" json:"priority"` // low | medium | high | urgent
	WorkspaceID  uint       `gorm:"not null;index" json:"workspace_id"`
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	DueDate      *time.Time `json:"due_date"`
	CompletedAt  *time.Time `json:"completed_at"`
	IsRecurring  bool       `gorm:"default:false" json:"is_recurring"`
	RecurRule    string     `json:"recur_rule"` // daily | weekly | monthly | weekdays
	ParentTaskID *uint      `json:"parent_task_id"` // for subtasks
	SortOrder    float64    `gorm:"default:0" json:"sort_order"` // LexoRank style
	Tags         []Tag      `gorm:"many2many:task_tags;constraint:OnDelete:CASCADE" json:"tags"`
	Comments     []Comment  `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Subtasks     []Task     `gorm:"foreignKey:ParentTaskID;constraint:OnDelete:CASCADE" json:"subtasks,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
