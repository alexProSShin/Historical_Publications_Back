package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func StartServer() {
	log.Println("Server start up")

	eventFile, err := os.Open("resources/data/events.json")
	if err != nil {
		log.Println("Ошибка при открытии JSON файла:", err)
		return
	}
	defer eventFile.Close()

	var events []HistoricalEvent
	decoder := json.NewDecoder(eventFile)
	if err := decoder.Decode(&events); err != nil {
		log.Println("Ошибка при декодировании JSON данных:", err)
		return
	}

	pubFile, err := os.Open("resources/data/publication.json")
	if err != nil {
		log.Println("Ошибка при открытии JSON файла публикаций:", err)
		return
	}
	defer pubFile.Close()

	var publication GetPublicationDTO
	pubDecoder := json.NewDecoder(pubFile)
	if err = pubDecoder.Decode(&publication); err != nil {
		log.Println("Ошибка при декодировании JSON данных публикаций:", err)
		return
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./resources/css")
	r.Static("/data", "./resources/data")
	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")

	r.GET("/", func(c *gin.Context) {
		title := c.DefaultQuery("title", "")
		var foundEvents []HistoricalEvent
		for _, event := range events {
			if strings.HasPrefix(strings.ToLower(event.Title), strings.ToLower(title)) {
				foundEvents = append(foundEvents, event)
			}
		}
		data := gin.H{
			"events":        foundEvents,
			"searchValue":   title,
			"publicationId": publication.ID,
			"eventsCount":   len(publication.Events),
		}
		c.HTML(http.StatusOK, "main_page.tmpl", data)
	})

	r.GET("/events/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		event := events[id-1]

		data := gin.H{
			"event":         event,
			"publicationId": publication.ID,
			"eventsCount":   len(publication.Events),
		}
		c.HTML(http.StatusOK, "event_page.tmpl", data)
	})

	r.GET("/publications/:id", func(c *gin.Context) {
		_, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		data := gin.H{
			"publication": publication,
			"eventsCount": len(publication.Events),
		}
		c.HTML(http.StatusOK, "publication_page.tmpl", data)
	})

	r.Run()

	log.Println("Server down")
}
