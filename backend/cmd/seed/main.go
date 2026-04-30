package main

import (
	"log"

	"github.com/GyBJluHv2/sensory-navigator-users/internal/auth"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/database"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/models"
	"gorm.io/gorm"
)

// Утилита заполняет БД демо-пользователем для ручного тестирования модуля.
func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("ошибка миграций: %v", err)
	}
	if err := database.SeedDemoPlaces(db); err != nil {
		log.Fatalf("ошибка инициализации мест: %v", err)
	}
	if err := seedDemoUser(db, cfg); err != nil {
		log.Fatalf("ошибка инициализации пользователя: %v", err)
	}
	log.Println("seed: готово, пользователь demo@example.com / demo123")
}

func seedDemoUser(db *gorm.DB, cfg *config.Config) error {
	var n int64
	db.Model(&models.User{}).Where("email = ?", "demo@example.com").Count(&n)
	if n > 0 {
		return nil
	}
	hash, err := auth.GeneratePasswordHash("demo123", cfg.BCryptCost)
	if err != nil {
		return err
	}
	user := &models.User{
		Email:        "demo@example.com",
		Username:     "demo",
		DisplayName:  "Демо пользователь",
		PasswordHash: hash,
		NoisePref:    3, LightPref: 3, CrowdPref: 3, SmellPref: 3, VisualPref: 3,
	}
	return db.Create(user).Error
}