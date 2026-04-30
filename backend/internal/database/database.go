package database

import (
	"errors"
	"log"

	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.PlaceRef{},
		&models.Review{},
		&models.Favorite{},
		&models.VerificationCode{},
	)
}

// SeedDemoPlaces заполняет таблицу places минимальными данными,
// чтобы продемонстрировать работу отзывов и избранного без модуля карты.
// При объединённой сборке seed выполняет модуль карты Атаханова Н. Р.
func SeedDemoPlaces(db *gorm.DB) error {
	demos := []models.PlaceRef{
		{ID: 1, Name: "Кафе \"Тихий дворик\""},
		{ID: 2, Name: "Парк Горького"},
		{ID: 3, Name: "Библиотека им. Достоевского"},
	}
	for _, p := range demos {
		var existing models.PlaceRef
		err := db.First(&existing, p.ID).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			if cerr := db.Create(&p).Error; cerr != nil {
				log.Printf("seed places: %v", cerr)
				return cerr
			}
		}
	}
	return nil
}