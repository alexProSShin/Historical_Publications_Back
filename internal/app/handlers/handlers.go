package handlers

import (
	"backend/internal/app/repository"
)

type Handler struct {
	repo  repository.PostgresRepo
	redis repository.RedisRepo
}

func New(repo repository.PostgresRepo, redis repository.RedisRepo) *Handler {
	return &Handler{
		repo:  repo,
		redis: redis,
	}
}
