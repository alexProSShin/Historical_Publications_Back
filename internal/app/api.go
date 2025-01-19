package app

import (
	"backend/internal/app/handlers"
	"backend/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func (app *Application) Run() {
	r := gin.Default()

	h := handlers.New(app.PostgresRepo, app.RedisRepo)
	app.setupRoutes(r, h)

	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server down")
}

func (app *Application) setupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.Use(middleware.InjectUserIDMiddleware())

	events := r.Group("/events")
	{
		events.GET("", h.HandleGetEvents)                                           // Получить все события
		events.GET("/:eventID", h.HandleGetEventByID)                               // Получить событие по ID
		events.POST("", h.HandleAddEvent)                                           // Добавить новое событие
		events.PUT("/:eventID", h.HandleUpdateEvent)                                // Обновить событие
		events.DELETE("/:eventID", h.HandleDeleteEvent)                             // Удалить событие
		events.POST("/:eventID", h.HandleAddEventImage)                             // Загрузить изображение для события
		events.POST("/:eventID/publications", h.HandleAddEventToPublication)        // Добавить событие в черновую публикацию
		events.DELETE("/:eventID/publications", h.HandleRemoveEventFromPublication) // Удалить событие из публикации
		events.PUT("/:eventID/publications", h.HandleUpdateEventPriority)           // Изменить приоритет события
	}

	publications := r.Group("/publications")
	{
		publications.GET("", h.HandleGetPublications)                               // Получить список публикаций с фильтрацией
		publications.GET("/:publicationID", h.HandleGetPublicationByID)             // Получить публикацию по ID
		publications.PUT("/:publicationID", h.HandleUpdatePublication)              // Обновить данные публикации
		publications.PUT("/:publicationID/form", h.HandleFormPublication)           // Сформировать публикацию
		publications.PUT("/:publicationID/finalized", h.HandleFinalizedPublication) // Одобрить/отменить публикацию
		publications.DELETE("/:publicationID", h.HandleDeletePublication)           // Удалить публикацию
	}

	users := r.Group("/users")
	{
		users.POST("/register", h.HandleRegisterUser) // Регистрация пользователя
		users.POST("/login", h.HandleLoginUser)       // Вход пользователя
		users.POST("/logout", h.HandleLogoutUser)     // Выход пользователя
		users.PUT("/me", h.HandleUpdateUser)          // Изменение данных пользователя
	}
}
