package models

import (
	"time"
)

type Workspace struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex:idx_user_slug" json:"slug"`
	Color     string    `gorm:"default:'#171717'" json:"color"`  // hex color
	Icon      string    `gorm:"default:'briefcase'" json:"icon"` // Tabler icon name
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_slug" json:"user_id"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsArchived bool     `gorm:"default:false" json:"is_archived"`
	Tasks     []Task    `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
