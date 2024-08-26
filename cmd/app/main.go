package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SavelyDev/crud-app/internal/config"
	"github.com/SavelyDev/crud-app/internal/repository/psql"
	"github.com/SavelyDev/crud-app/internal/service"
	"github.com/SavelyDev/crud-app/internal/transport/rest"
	"github.com/SavelyDev/crud-app/pkg/database"
	"github.com/SavelyDev/crud-app/pkg/hash"
	"github.com/SavelyDev/crud-app/pkg/server"
)

// @title CRUD-APP API
// @version 1.0

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

const (
	CONFIG_DIR  = "../../configs"
	CONFIG_FILE = "config"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := database.New(database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher(cfg.Hash.Slat)

	authRepo := psql.NewAuthRepo(db)
	todoListRepo := psql.NewTodoListRepo(db)
	todoItemRepo := psql.NewTodoItemRepo(db)

	authService := service.NewAuthService(authRepo, hasher, cfg.Auth.TokenTTL, cfg.Auth.Secret)
	todoListService := service.NewTodoListService(todoListRepo)
	todoItemService := service.NewTodoItemService(todoItemRepo, todoListRepo)

	hand := rest.NewHandler(authService, todoListService, todoItemService)

	srv := server.NewServer(cfg.Server.Port, hand.InitRouter())
	go func() {
		if err := srv.Run(); err != nil {
			logrus.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
}
