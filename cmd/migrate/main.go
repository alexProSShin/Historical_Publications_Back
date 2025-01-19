package main

import (
	"backend/internal/app/dsn"
	"backend/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
			Title:       "Первый полет в космос",
			EventType:   models.EventType("Событие"),
			Description: "12 апреля 1961 года Юрий Гагарин стал первым человеком, полетевшим в космос на корабле 'Восток 1'.",
			Info:        "Первый космический полет Юрия Гагарина.",
			PhotoURL:    stringPointer("http://localhost:9000/events/1.png"),
			Source:      stringPointer("http://ru.wikipedia.org/wiki/Гагарин,_Юрий_Алексеевич"),
		},
		{
			Title:       "Первая мировая война",
			EventType:   models.EventType("Событие"),
			Description: "Первая мировая война (1914-1918) была одним из самых разрушительных конфликтов в истории человечества.",
			Info:        "Глобальный конфликт, который повлек за собой гибель миллионов людей и изменение политической карты мира.",
			PhotoURL:    stringPointer("http://localhost:9000/events/2.png"),
			Source:      stringPointer("http://ru.wikipedia.org/wiki/Первая_мировая_война"),
		},
		{
			Title:       "Революция 1917 года в России",
			EventType:   models.EventType("Событие"),
			Description: "В 1917 году в России произошла Октябрьская революция, приведшая к свержению царизма и установлению советской власти.",
			Info:        "Октябрьская революция стала важнейшим событием в истории России и мира, приведя к созданию Советского Союза.",
			PhotoURL:    stringPointer("http://localhost:9000/events/3.png"),
			Source:      stringPointer("http://ru.wikipedia.org/wiki/Октябрьская_революция"),
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
