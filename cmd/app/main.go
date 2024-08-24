package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SavelyDev/crud-app/internal/repository/psql"
	"github.com/SavelyDev/crud-app/internal/service"
	"github.com/SavelyDev/crud-app/internal/transport/rest"
	"github.com/SavelyDev/crud-app/pkg/database"
	"github.com/SavelyDev/crud-app/pkg/server"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// @title CRUD-APP API
// @version 1.0

// @host localhost:80
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err := gotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	db, err := database.NewDB(database.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	authRepo := psql.NewAuthRepo(db)
	todoListRepo := psql.NewTodoListRepo(db)
	todoItemRepo := psql.NewTodoItemRepo(db)

	authService := service.NewAuthService(authRepo)
	todoListService := service.NewTodoListService(todoListRepo)
	todoItemService := service.NewTodoItemService(todoItemRepo, todoListRepo)

	hand := rest.NewHandler(authService, todoListService, todoItemService)

	srv := server.NewServer(viper.GetString("server.port"), hand.InitRouter())
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

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
