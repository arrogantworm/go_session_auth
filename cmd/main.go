package main

import (
	"context"
	"log"
	"os"
	"session-auth/handler"
	"session-auth/repository"
	"session-auth/server"
	"session-auth/service"
	"time"

	_ "session-auth/cmd/docs"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// @title Session Auth
// @version 1.0
// @description Session auth project

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name session_token

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	cfg := repository.Config{
		User:     os.Getenv("PGUser"),
		Password: os.Getenv("PGPassword"),
		Host:     os.Getenv("PGHost"),
		Port:     os.Getenv("PGPort"),
		DBName:   os.Getenv("PGDBName"),
		SSLMode:  os.Getenv("PGSSLMode"),
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	dbctx, cancel := context.WithTimeout(ctx, time.Second*10)
	database, err := repository.NewPostgres(dbctx, &cfg)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewRepository(database)
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	server := server.NewServer("8000", handler, time.Second*10)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
