package http

import (
	"net/http"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTaskRequestBody struct {
	Title string `json:"title"`
}

type CreateTaskResponseData struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

func CreateTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req CreateTaskRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := app.taskService.CreateTask(gtx, task.CreateTaskRequest{Title: req.Title})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, CreateTaskResponseData{
			ID:    t.ID,
			Title: t.Title,
		})
	}
}
