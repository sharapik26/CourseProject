package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/database"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/routes"
)

// @title       Sensory Navigator — Users & Reviews API
// @version     1.0
// @description REST API модуля пользователей и отзывов проекта «Сенсорный навигатор».
// @host        localhost:8081
// @BasePath    /api
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, используются переменные окружения")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("ошибка миграций схемы БД: %v", err)
	}

	if err := database.SeedDemoPlaces(db); err != nil {
		log.Fatalf("ошибка инициализации демо-мест: %v", err)
	}

	router := routes.NewRouter(db, cfg)

	addr := ":" + cfg.Port
	log.Printf("запуск сервера на %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("ошибка запуска HTTP-сервера: %v", err)
		os.Exit(1)
	}
}