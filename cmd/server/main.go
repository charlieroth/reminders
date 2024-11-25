package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/charlieroth/reminders/internal/config"
	remindersHttp "github.com/charlieroth/reminders/internal/http"
	"github.com/charlieroth/reminders/internal/outbound"
	"github.com/charlieroth/reminders/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
	)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// catching ctx.Done() timeout of 5 seconds
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
