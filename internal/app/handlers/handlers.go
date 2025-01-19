package handlers

import (
	"backend/internal/app/repository"
)

// @title           RIP API
// @version         1.0
// @description     API для управления событиями, публикациями и изображениями.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @schemes   http

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization
// @description Use the JWT Bearer Token for authorization. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

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
