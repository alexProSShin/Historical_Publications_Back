package app

import (
	"backend/internal/app/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func (app *Application) Run() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./resources/css")
	r.Static("/data", "./resources/data")
	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")

	h := handlers.New(app.Repository)
	app.setupRoutes(r, h)

	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server down")
}

func (app *Application) setupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.GET("/", h.HandleGetEvents)
	r.GET("/events/:id", h.HandleGetEventByID)
	r.POST("/events/:eventID/publications", h.HandleAddEventToPublication)
	r.GET("/publications/:id", h.HandleGetPublicationByID)
}
