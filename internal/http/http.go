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
	userService     *service.UserService
	databaseService *service.DatabaseService
	taskService     *service.TaskService
	listService     *service.ListService
}

func NewHttpServer(
	userService *service.UserService,
	databaseService *service.DatabaseService,
	taskService *service.TaskService,
	listService *service.ListService,
	config HttpServerConfig,
) *http.Server {
	app := &App{userService: userService, databaseService: databaseService, taskService: taskService, listService: listService}
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

	router.GET("/users", GetUsers(app))
	router.POST("/users", CreateUser(app))
	router.GET("/users/:user_id", GetUser(app))
	router.PATCH("/users/:user_id", UpdateUser(app))

	router.GET("/lists", GetLists(app))
	router.GET("/lists/:list_id", GetList(app))
	router.POST("/lists", CreateList(app))
	router.PATCH("/lists/:list_id", UpdateList(app))

	router.GET("/lists/:list_id/tasks", GetListTasks(app))
	router.POST("/lists/:list_id/tasks", CreateListTask(app))
	router.GET("/lists/:list_id/tasks/:task_id", GetListTask(app))
	router.PATCH("/lists/:list_id/tasks/:task_id", UpdateListTask(app))
}
