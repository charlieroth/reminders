package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ReadinessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func LivenessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

type CreateTaskRequestBody struct {
	Title string `json:"title"`
}

type CreateTaskResponseData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateTaskResponseData(task task.Task) CreateTaskResponseData {
	return CreateTaskResponseData{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

type GetTaskResponseData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetTaskResponseData(task task.Task) GetTaskResponseData {
	return GetTaskResponseData{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func NewListTasksResponseData(tasks []task.Task) []GetTaskResponseData {
	responseData := []GetTaskResponseData{}
	for _, t := range tasks {
		responseData = append(responseData, NewGetTaskResponseData(t))
	}
	return responseData
}

func CreateTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req CreateTaskRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Title == "" {
			err := &task.TaskTitleEmptyError{}
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := app.taskService.CreateTask(gtx, task.NewCreateTaskRequest(req.Title))
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusCreated, NewCreateTaskResponseData(t))
	}
}

func ListTasks(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		tasks, err := app.taskService.ListTasks(gtx)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewListTasksResponseData(tasks))
	}
}

func GetTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		id, err := uuid.Parse(gtx.Param("id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := app.taskService.GetTask(gtx, id)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetTaskResponseData(t))
	}
}
