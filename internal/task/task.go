package task

import "github.com/google/uuid"

type Task struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}
