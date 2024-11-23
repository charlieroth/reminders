package http

import (
	"fmt"
	"net/http"

	"github.com/charlieroth/reminders/internal/service"
	"github.com/gin-gonic/gin"
)

type HttpServerConfig struct {
	Port string
}

type App struct {
	taskService *service.TaskService
	listService *service.ListService
}

func NewHttpServer(taskService *service.TaskService, listService *service.ListService, config HttpServerConfig) *http.Server {
	app := &App{taskService: taskService, listService: listService}
	router := gin.New()
	apiRoutes(router, app)
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router.Handler(),
	}
}

func apiRoutes(router *gin.Engine, app *App) {
	router.GET("/openapi.json", ServeOpenAPISpec(app))

	router.GET("/readiness", ReadinessCheck(app))
	router.GET("/liveness", LivenessCheck(app))

	router.GET("/tasks", ListTasks(app))
	router.POST("/tasks", CreateTask(app))
	router.GET("/tasks/:id", GetTask(app))
	router.PATCH("/tasks/:id", UpdateTask(app))

	router.GET("/lists", ListLists(app))
	router.GET("/lists/:id", GetList(app))
	router.POST("/lists", CreateList(app))
	router.PATCH("/lists/:id", UpdateList(app))
}
