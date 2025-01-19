package main

import (
	"backend/internal/app/config"
	"backend/internal/app/dsn"
	postgresrepo "backend/internal/app/repository/postgres"
	redisrepo "backend/internal/app/repository/redis"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
		return
	}

	postgresRepo := postgresrepo.NewPostgresRepository(db, minioClient, cfg.MinioBucket)
	if err != nil {
		log.Fatal(err)
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
		return
	}

	// Создаем репозиторий Redis
	redisRepo := redisrepo.NewRedisRepository(redisClient)

	application, err := app.New(cfg, postgresRepo, redisRepo)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	application.Run()
	log.Println("Application terminated!")
}
