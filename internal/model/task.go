package model

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	AuthorID    uint      `gorm:"not null" json:"-"`
	Author      User      `gorm:"foreignKey:AuthorID" json:"author"`
	ExecutorID  uint      `gorm:"not null" json:"-"`
	Executor    User      `gorm:"foreignKey:ExecutorID" json:"executor"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDone      bool      `gorm:"default:false" json:"is_done"`
}
