package http

import (
	"fmt"
	"net/http"

	"github.com/charlieroth/reminders/internal/config"
	"github.com/charlieroth/reminders/internal/service"
	"github.com/gin-gonic/gin"
)

type HttpServerConfig struct {
	Port string
}

type App struct {
	userService     *service.UserService
	authService     *service.AuthService
	databaseService *service.DatabaseService
	taskService     *service.TaskService
	listService     *service.ListService
	config          *config.Config
}

func NewHttpServer(
	userService *service.UserService,
	authService *service.AuthService,
	databaseService *service.DatabaseService,
	taskService *service.TaskService,
	listService *service.ListService,
	config *config.Config,
) *http.Server {
	app := &App{
		config:          config,
		userService:     userService,
		authService:     authService,
		databaseService: databaseService,
		taskService:     taskService,
		listService:     listService,
	}
	router := gin.New()
	apiRoutes(router, app)
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.ServerPort),
		Handler: router.Handler(),
	}
}

func apiRoutes(router *gin.Engine, app *App) {
	router.GET("/openapi.json", ServeOpenAPISpec(app))

	router.GET("/readiness", ReadinessCheck(app))
	router.GET("/liveness", LivenessCheck(app))

	authGroup := router.Group("/auth")
	authGroup.POST("/login", Login(app))
	authGroup.POST("/logout", Logout(app))
	authGroup.POST("/refresh", Refresh(app))
	authGroup.GET("/sessions", GetSessions(app))

	userGroup := router.Group("/users")
	userGroup.GET("", GetUsers(app))
	userGroup.POST("", CreateUser(app))
	userGroup.GET("/:user_id", GetUser(app))
	userGroup.PATCH("/:user_id", UpdateUser(app))

	listGroup := router.Group("/lists")
	listGroup.Use(AuthMiddleware(app))
	listGroup.GET("", GetLists(app))
	listGroup.GET("/:list_id", GetList(app))
	listGroup.POST("", CreateList(app))
	listGroup.PATCH("/:list_id", UpdateList(app))

	taskGroup := listGroup.Group("/:list_id/tasks")
	taskGroup.Use(AuthMiddleware(app))
	taskGroup.GET("", GetListTasks(app))
	taskGroup.POST("", CreateListTask(app))
	taskGroup.GET("/:task_id", GetListTask(app))
	taskGroup.PATCH("/:task_id", UpdateListTask(app))
}
