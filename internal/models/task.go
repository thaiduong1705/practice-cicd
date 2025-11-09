package models

// Task represents a TODO item. The ID is a string primary key so clients can
// supply e.g. a UUID. Completed stored as boolean; GORM maps it to INTEGER.
type Task struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"not null" validate:"required"`
	Description string `json:"description" validate:"required"`
	Completed   bool   `json:"completed" gorm:"not null;default:false"`
}
