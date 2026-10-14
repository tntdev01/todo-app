package main

// @title Todo API
// @version 1.0
// @description This is a simple Todo API built with Go and Gin.
// @host localhost:8080
// @BasePath /

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tuanntran/todo-app/todo-api/docs"

	"github.com/tuanntran/todo-app/todo-api/internal/config"
	"github.com/tuanntran/todo-app/todo-api/internal/database"
	"github.com/tuanntran/todo-app/todo-api/internal/handler"
	"github.com/tuanntran/todo-app/todo-api/internal/repository"
	"github.com/tuanntran/todo-app/todo-api/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("application starting")

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Info("configuration loaded", "port", cfg.AppPort)

	conn, err := database.Connect(
		cfg.DatabaseURL(),
	)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		panic(err)
	}
	logger.Info("database connected")

	repo := repository.NewTodoRepository(conn)

	svc := service.NewTodoService(repo)

	h := &handler.TodoHandler{
		Service: svc,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/todos", h.Create)
	router.GET("/todos", h.GetAll)
	router.GET("/todos/:id", h.GetByID)
	router.PUT("/todos/:id", h.Update)
	router.DELETE("/todos/:id", h.Delete)

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// router.Run(":" + cfg.AppPort)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		logger.Info("http server starting", "port", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start http server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("http server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shut down http server", "error", err)
	}

	logger.Info("closing database connections")
	conn.Close()
	logger.Info("application stopped")
}
