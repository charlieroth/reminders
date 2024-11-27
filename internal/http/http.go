package http

import (
	"fmt"
	"net/http"

	"github.com/charlieroth/reminders/internal/config"
	"github.com/charlieroth/reminders/internal/service"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	logger          *zerolog.Logger
}

func NewHttpServer(
	userService *service.UserService,
	authService *service.AuthService,
	databaseService *service.DatabaseService,
	taskService *service.TaskService,
	listService *service.ListService,
	config *config.Config,
	logger *zerolog.Logger,
) *http.Server {
	app := &App{
		logger:          logger,
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
	router.Use(requestid.New())
	router.Use(logger.SetLogger())

	// GET:/openapi.json
	router.GET("/openapi.json", ServeOpenAPISpec(app))
	// GET:/readiness
	router.GET("/readiness", ReadinessCheck(app))
	// GET:/liveness
	router.GET("/liveness", LivenessCheck(app))

	adminGroup := router.Group("/admin")
	{
		adminGroup.Use(AdminMiddleware(app))
		// GET:/admin/users
		adminGroup.GET("/users", GetUsers(app))
		// GET:/admin/users/:user_id
		adminGroup.GET("/users/:user_id", GetUser(app))
	}

	listGroup := router.Group("/lists")
	{
		listGroup.Use(AuthMiddleware(app))
		// GET:/lists
		listGroup.GET("", GetLists(app))
		// GET:/lists/:list_id
		listGroup.GET("/:list_id", GetList(app))
		// POST:/lists
		listGroup.POST("", CreateList(app))
		// PATCH:/lists/:list_id
		listGroup.PATCH("/:list_id", UpdateList(app))
	}

	taskGroup := listGroup.Group("/:list_id/tasks")
	{
		taskGroup.Use(AuthMiddleware(app))
		// GET:/lists/:list_id/tasks
		taskGroup.GET("", GetListTasks(app))
		// POST:/lists/:list_id/tasks
		taskGroup.POST("", CreateListTask(app))
		// GET:/lists/:list_id/tasks/:task_id
		taskGroup.GET("/:task_id", GetListTask(app))
		// PATCH:/lists/:list_id/tasks/:task_id
		taskGroup.PATCH("/:task_id", UpdateListTask(app))
	}

	userGroup := router.Group("/users")
	// POST:/users
	userGroup.POST("", Register(app))
	// POST:/users/login
	userGroup.POST("/login", Login(app))
	{
		userAuthGroup := userGroup.Group("/")
		userAuthGroup.Use(AuthMiddleware(app))
		{
			// PATCH:/users/:user_id
			userAuthGroup.PATCH("/:user_id", UpdateUser(app))
			// POST:/users/logout
			userAuthGroup.POST("/logout", Logout(app))
		}
	}

	tokenGroup := router.Group("/tokens")
	{
		tokenGroup.Use(AuthMiddleware(app))
		// POST:/tokens/refresh
		tokenGroup.POST("/refresh", Refresh(app))
		// POST:/tokens/revoke/:session_id
		tokenGroup.POST("/revoke/:session_id", RevokeSession(app))

	}
}
