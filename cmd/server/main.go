package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/charlieroth/reminders/internal/config"
	remindersHttp "github.com/charlieroth/reminders/internal/http"
	"github.com/charlieroth/reminders/internal/logger"
	"github.com/charlieroth/reminders/internal/outbound"
	"github.com/charlieroth/reminders/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	config, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := logger.NewLogger(config.Environment)

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		logger.Error().Msgf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	pg := outbound.NewPg(db)
	userService := service.NewUserService(pg)
	authService := service.NewAuthService(pg)
	databaseService := service.NewDatabaseService(pg)
	taskService := service.NewTaskService(pg)
	listService := service.NewListService(pg)
	srv := remindersHttp.NewHttpServer(
		userService,
		authService,
		databaseService,
		taskService,
		listService,
		config,
		logger,
	)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Msgf("listen: %s\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Msgf("Server forced to shutdown: %v\n", err)
	}

	// catching ctx.Done() timeout of 5 seconds
	select {
	case <-ctx.Done():
		logger.Error().Msg("timeout of 5 seconds.")
	}
	logger.Info().Msg("Server exiting")
}
