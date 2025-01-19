package app

import (
	"backend/internal/app/config"
	"backend/internal/app/repository"
)

type Application struct {
	Config       *config.Config
	PostgresRepo repository.PostgresRepo
	RedisRepo    repository.RedisRepo
}

func New(cfg *config.Config, postgresRepo repository.PostgresRepo, redisRepo repository.RedisRepo) (*Application, error) {
	return &Application{
		Config:       cfg,
		PostgresRepo: postgresRepo,
		RedisRepo:    redisRepo,
	}, nil
}
