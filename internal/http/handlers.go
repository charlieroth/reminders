package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/list"
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

func NewListTasksResponseData(tasks []task.Task) []GetTaskResponseData {
	responseData := []GetTaskResponseData{}
	for _, t := range tasks {
		responseData = append(responseData, NewGetTaskResponseData(t))
	}
	return responseData
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

type CreateTaskRequestBody struct {
	Title string `json:"title"`
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

type UpdateTaskRequestBody struct {
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

func UpdateTask(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req UpdateTaskRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := uuid.Parse(gtx.Param("id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := app.taskService.UpdateTask(gtx, id, task.NewUpdateTaskRequest(req.Title, req.Completed))
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewUpdateTaskResponseData(t))
	}
}

type CreateListRequestBody struct {
	Name string `json:"name"`
}

type CreateListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateListResponseData(list list.List) CreateListResponseData {
	return CreateListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func CreateList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req CreateListRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.CreateList(gtx, list.NewCreateListRequest(req.Name))
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusCreated, NewCreateListResponseData(l))
	}
}

type UpdateListRequestBody struct {
	Name string `json:"name"`
}

func NewUpdateListRequest(name string) list.UpdateListRequest {
	return list.UpdateListRequest{Name: name}
}

type UpdateListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUpdateListResponseData(list list.List) UpdateListResponseData {
	return UpdateListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func UpdateList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req UpdateListRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := uuid.Parse(gtx.Param("id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.UpdateList(gtx, id, NewUpdateListRequest(req.Name))
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewUpdateListResponseData(l))
	}
}

type GetListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetListResponseData(list list.List) GetListResponseData {
	return GetListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func GetList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		id, err := uuid.Parse(gtx.Param("id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.GetList(gtx, id)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetListResponseData(l))
	}
}

func NewListListsResponseData(lists []list.List) []GetListResponseData {
	responseData := []GetListResponseData{}
	for _, l := range lists {
		responseData = append(responseData, NewGetListResponseData(l))
	}
	return responseData
}

func ListLists(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		lists, err := app.listService.GetLists(gtx)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewListListsResponseData(lists))
	}
}
