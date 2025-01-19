package main

import (
	"backend/internal/app/config"
	"backend/internal/app/dsn"
	"backend/internal/app/repository"
	"context"
	"log"

	"backend/internal/app"
)

func main() {
	log.Println("Application start!")
	ctx := context.Background()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	repo, err := repository.NewPostgresRepository(dsn.FromEnv())
	if err != nil {
		log.Fatal(err)
		return
	}

	application, err := app.New(cfg, repo)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	application.Run()
	log.Println("Application terminated!")
}
