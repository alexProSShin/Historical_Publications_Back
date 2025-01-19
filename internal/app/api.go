package app

import (
	_ "backend/docs"
	"backend/internal/app/handlers"
	"backend/internal/middleware"
	"backend/internal/middleware/cors"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func (app *Application) Run() {
	r := gin.Default()

	h := handlers.New(app.PostgresRepo, app.RedisRepo)
	app.setupRoutes(r, h)

	r.Use(cors.New(cors.CORSOptions{
		AllowedOrigins:   []string{"http://localhost:3000", "https://your-production-domain.com"}, // Укажите конкретные домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Disposition", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server down")
}

func (app *Application) setupRoutes(r *gin.Engine, h *handlers.Handler) {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	events := r.Group("/events")
	{
		events.Use(middleware.AuthMiddleware(app.RedisRepo))
		events.POST("", middleware.RequireModeratorRole(h.GetUserByID), h.HandleAddEvent)
		events.PUT("/:eventID", middleware.RequireModeratorRole(h.GetUserByID), h.HandleUpdateEvent)
		events.DELETE("/:eventID", middleware.RequireModeratorRole(h.GetUserByID), h.HandleDeleteEvent)
		events.POST("/:eventID", middleware.RequireModeratorRole(h.GetUserByID), h.HandleAddEventImage) // Загрузить изображение для события
		events.POST("/:eventID/publications", h.HandleAddEventToPublication)                            // Добавить событие в черновую публикацию
		events.DELETE("/:eventID/publications", h.HandleRemoveEventFromPublication)                     // Удалить событие из публикации
		events.PUT("/:eventID/publications", h.HandleUpdateEventPriority)                               // Изменить приоритет события

	}

	publications := r.Group("/publications")
	{
		publications.Use(middleware.AuthMiddleware(app.RedisRepo))
		publications.GET("", h.HandleGetPublications)                                                                               // Получить список публикаций с фильтрацией
		publications.GET("/:publicationID", h.HandleGetPublicationByID)                                                             // Получить публикацию по ID
		publications.PUT("/:publicationID", h.HandleUpdatePublication)                                                              // Обновить данные публикации
		publications.PUT("/:publicationID/form", h.HandleFormPublication)                                                           // Сформировать публикацию
		publications.PUT("/:publicationID/finalized", middleware.RequireModeratorRole(h.GetUserByID), h.HandleFinalizedPublication) // Одобрить/отменить публикацию
		publications.DELETE("/:publicationID", h.HandleDeletePublication)                                                           // Удалить публикацию
	}

	users := r.Group("/users")
	{
		users.Use(middleware.AuthMiddleware(app.RedisRepo))
		users.PUT("/me", h.HandleUpdateUser)
		users.POST("/logout", h.HandleLogoutUser)
	}

	users = r.Group("/users")
	{
		users.POST("/register", h.HandleRegisterUser)
		users.POST("/login", h.HandleLoginUser)
	}

	events = r.Group("/events")
	{
		events.GET("", h.HandleGetEvents)
		events.GET("/:eventID", h.HandleGetEventByID)
	}
}
