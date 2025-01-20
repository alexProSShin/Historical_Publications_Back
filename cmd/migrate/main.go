package main

import (
	"backend/internal/app/dsn"
	"backend/internal/models"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.HistoricalEvent{}, &models.PublicationsEvents{}, &models.Publication{}, &models.User{})
	if err != nil {
		panic("cant migrate db")
	}

	err = insertTestData(db)
	if err != nil {
		log.Printf("Ошибка вставки тестовых данных: %v", err)
	} else {
		log.Println("Тестовые данные успешно вставлены.")
	}
}

func insertTestData(db *gorm.DB) error {
	testEvents := []models.HistoricalEvent{
		{
			Title:       "Колизей в Риме",
			EventType:   models.EventType("Локация"),
			Description: "Знаменитый амфитеатр в Риме.",
			Info:        "Колизей, построенный в 70-80 годах нашей эры, является одним из самых известных памятников архитектуры в мире. Он использовался для гладиаторских боев, театральных представлений и других массовых мероприятий. Подробнее: https://example.com/colosseum",
			PhotoURL:    stringPointer("http://localhost:9000/events/1.png"),
			Source:      stringPointer("Исторические хроники"),
		},
		{
			Title:       "Великие пирамиды Гизы",
			EventType:   models.EventType("Локация"),
			Description: "Комплекс пирамид в Египте.",
			Info:        "Великие пирамиды Гизы – это древнейший и единственный сохранившийся памятник из семи чудес света. Строительство пирамид датируется примерно 2600 годом до н. э. Подробнее: https://example.com/pyramids",
			PhotoURL:    stringPointer("http://localhost:9000/events/2.jpg"),
			Source:      stringPointer("Историческая энциклопедия"),
		},
		{
			Title:       "Битва при Ватерлоо",
			EventType:   models.EventType("Событие"),
			Description: "Последнее сражение Наполеона.",
			Info:        "Битва при Ватерлоо, состоявшаяся 18 июня 1815 года, стала последним крупным сражением Наполеона Бонапарта. Его поражение ознаменовало конец эпохи наполеоновских войн. Подробнее: https://example.com/waterloo",
			PhotoURL:    stringPointer("http://localhost:9000/events/3.jpg"),
			Source:      stringPointer("Военная история"),
		},
		{
			Title:       "Высадка на Луну",
			EventType:   models.EventType("Событие"),
			Description: "Первый шаг человека на Луну.",
			Info:        "20 июля 1969 года Нил Армстронг стал первым человеком, ступившим на поверхность Луны в рамках миссии Apollo 11. Эта миссия изменила ход истории освоения космоса. Подробнее: https://example.com/moon_landing",
			PhotoURL:    stringPointer("http://localhost:9000/events/4.webp"),
			Source:      stringPointer("Космическая энциклопедия"),
		},
		{
			Title:       "Скрижали Моисея",
			EventType:   models.EventType("Артефакт"),
			Description: "Священные каменные таблички.",
			Info:        "Скрижали Моисея – это каменные таблички, на которых, согласно Библии, были начертаны Десять заповедей, данные Моисею на горе Синай. Подробнее: https://example.com/ten_commandments",
			PhotoURL:    stringPointer("http://localhost:9000/events/5.jpg"),
			Source:      stringPointer("Религиозные тексты"),
		},
		{
			Title:       "Скипетр Тутанхамона",
			EventType:   models.EventType("Артефакт"),
			Description: "Древний египетский артефакт.",
			Info:        "Скипетр Тутанхамона – один из множества артефактов, найденных в гробнице этого фараона. Этот символ власти был создан более 3000 лет назад. Подробнее: https://example.com/tutankhamun_scepter",
			PhotoURL:    stringPointer("http://localhost:9000/events/6.jpg"),
			Source:      stringPointer("Археологические находки"),
		},
	}

	if err := db.Create(&testEvents).Error; err != nil {
		return err
	}

	return nil
}

func stringPointer(s string) *string {
	return &s
}
