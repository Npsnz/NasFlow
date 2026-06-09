package models

type Tag struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID uint   `gorm:"index" json:"user_id"`
	Tasks  []Task `gorm:"many2many:task_tags;" json:"-"` // Hidden in JSON to avoid circular references
}
