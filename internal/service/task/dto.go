package task

import (
	"mine/internal/model"
	"time"
)

// --- Request DTOs ---

type CreateTaskRequest struct {
	Name        string    `json:"name" validate:"required,min=3"`
	Description string    `json:"description" validate:"max=500"`
	AuthorID    uint      `json:"author_id" validate:"required"`
	ExecutorID  uint      `json:"executor_id" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
}

type UpdateTaskRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=3"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	IsDone      *bool      `json:"is_done,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

// --- Response DTOs ---

type TaskResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Author      TaskUserResponse `json:"author"`
	Executor    TaskUserResponse `json:"executor"`
	Deadline    time.Time        `json:"deadline"`
	CreatedAt   time.Time        `json:"created_at"`
	IsDone      bool             `json:"is_done"`
}

type TaskUserResponse struct {
	ID               uint   `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	TelegramUsername string `json:"telegram_username" validate:"required"`
}

// --- Helpers ---

func ToTaskResponse(t model.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Author: TaskUserResponse{
			ID:               t.Author.ID,
			Name:             t.Author.Name,
			TelegramUsername: t.Author.TelegramUsername,
		},
		Executor: TaskUserResponse{
			ID:               t.Executor.ID,
			Name:             t.Executor.Name,
			TelegramUsername: t.Executor.TelegramUsername,
		},
		Deadline:  t.Deadline,
		CreatedAt: t.CreatedAt,
		IsDone:    t.IsDone,
	}
}
