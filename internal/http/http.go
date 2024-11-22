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
	taskService *service.Service
}

func NewHttpServer(taskService *service.Service, config HttpServerConfig) *http.Server {
	app := &App{taskService: taskService}
	router := gin.New()
	apiRoutes(router, app)
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router.Handler(),
	}
}

func apiRoutes(router *gin.Engine, app *App) {
	router.POST("/tasks", CreateTask(app))
	// router.PATCH("/tasks/:task_id", UpdateTask(app))
	// router.GET("/tasks/:task_id", GetTask(app))
}
