package app

import (
	"backend/internal/app/config"
	"backend/internal/app/repository"
)

type Application struct {
	Config     *config.Config
	Repository repository.Repository
}

func New(cfg *config.Config, repo repository.Repository) (*Application, error) {
	return &Application{
		Config:     cfg,
		Repository: repo,
	}, nil
}
