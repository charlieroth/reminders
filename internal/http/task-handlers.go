package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetTaskResponseData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetTaskResponseData(task task.Task) GetTaskResponseData {
	return GetTaskResponseData{
		ID:        task.ID,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

type GetListTaskResponseData struct {
	ListID uuid.UUID           `json:"list_id"`
	Task   GetTaskResponseData `json:"task"`
}

func NewGetListTaskResponseData(task task.Task) GetListTaskResponseData {
	return GetListTaskResponseData{
		Task: NewGetTaskResponseData(task),
	}
}

func GetListTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		taskID, err := uuid.Parse(gtx.Param("task_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := app.taskService.GetListTask(gtx, listID, taskID)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetListTaskResponseData(t))
	}
}

type GetListTasksResponseData struct {
	Tasks []GetListTaskResponseData `json:"tasks"`
}

func NewGetListTasksResponseData(tasks []task.Task) GetListTasksResponseData {
	responseData := []GetListTaskResponseData{}
	for _, t := range tasks {
		responseData = append(responseData, NewGetListTaskResponseData(t))
	}

	return GetListTasksResponseData{
		Tasks: responseData,
	}
}

func GetListTasks(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tasks, err := app.taskService.GetListTasks(gtx, listID)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetListTasksResponseData(tasks))
	}
}

type CreateListTaskRequest struct {
	Title string `json:"title"`
}

func NewCreateListTaskRequest(gtx *gin.Context) (CreateListTaskRequest, error) {
	var request CreateListTaskRequest
	if err := gtx.ShouldBindJSON(&request); err != nil {
		return CreateListTaskRequest{}, err
	}
	return request, nil
}

type CreateTaskResponseData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateTaskResponseData(task task.Task) CreateTaskResponseData {
	return CreateTaskResponseData{
		ID:        task.ID,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

type CreateListTaskResponseData struct {
	ListID uuid.UUID              `json:"list_id"`
	Task   CreateTaskResponseData `json:"task"`
}

func NewCreateListTaskResponseData(listID uuid.UUID, task task.Task) CreateListTaskResponseData {
	return CreateListTaskResponseData{
		ListID: listID,
		Task:   NewCreateTaskResponseData(task),
	}
}

func CreateListTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		request, err := NewCreateListTaskRequest(gtx)
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createTaskRequest := task.CreateTaskRequest{Title: request.Title}
		t, err := app.taskService.CreateListTask(gtx, listID, createTaskRequest)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewCreateListTaskResponseData(listID, t))
	}
}

type UpdateListTaskRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type UpdateTaskResponseData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUpdateTaskResponseData(task task.Task) UpdateTaskResponseData {
	return UpdateTaskResponseData{
		ID:        task.ID,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

type UpdateListTaskResponseData struct {
	ListID uuid.UUID              `json:"list_id"`
	Task   UpdateTaskResponseData `json:"task"`
}

func NewUpdateListTaskResponseData(listID uuid.UUID, task task.Task) UpdateListTaskResponseData {
	return UpdateListTaskResponseData{
		ListID: listID,
		Task:   NewUpdateTaskResponseData(task),
	}
}

func UpdateListTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var request UpdateListTaskRequest
		if err := gtx.ShouldBindJSON(&request); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		taskID, err := uuid.Parse(gtx.Param("task_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateListTaskRequest := task.UpdateTaskRequest{
			Title:     request.Title,
			Completed: request.Completed,
		}
		t, err := app.taskService.UpdateListTask(gtx, listID, taskID, updateListTaskRequest)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewUpdateListTaskResponseData(listID, t))
	}
}
