package services

import (
	"errors"

	"github.com/GyBJluHv2/sensory-navigator-users/internal/models"
	"gorm.io/gorm"
)

var (
	ErrReviewExists = errors.New("отзыв этого пользователя для места уже существует")
	ErrNotOwnReview = errors.New("отзыв принадлежит другому пользователю")
	ErrPlaceMissing = errors.New("место не найдено")
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) Create(r *models.Review) error {
	var place models.PlaceRef
	if err := s.db.First(&place, r.PlaceID).Error; err != nil {
		return ErrPlaceMissing
	}
	var n int64
	s.db.Model(&models.Review{}).
		Where("user_id = ? AND place_id = ?", r.UserID, r.PlaceID).
		Count(&n)
	if n > 0 {
		return ErrReviewExists
	}
	return s.db.Create(r).Error
}

func (s *ReviewService) Update(userID uint64, r *models.Review) error {
	var existing models.Review
	if err := s.db.First(&existing, r.ID).Error; err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrNotOwnReview
	}
	existing.Text = r.Text
	existing.Noise = r.Noise
	existing.Light = r.Light
	existing.Crowd = r.Crowd
	existing.Smell = r.Smell
	existing.Visual = r.Visual
	return s.db.Save(&existing).Error
}

func (s *ReviewService) Delete(userID, reviewID uint64) error {
	var existing models.Review
	if err := s.db.First(&existing, reviewID).Error; err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrNotOwnReview
	}
	return s.db.Delete(&existing).Error
}

func (s *ReviewService) ListByPlace(placeID uint64) ([]models.Review, error) {
	var rs []models.Review
	err := s.db.Preload("User").
		Where("place_id = ?", placeID).
		Order("created_at DESC").
		Find(&rs).Error
	return rs, err
}

func (s *ReviewService) ListByUser(userID uint64) ([]models.Review, error) {
	var rs []models.Review
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&rs).Error
	return rs, err
}

// Aggregate возвращает агрегированные средние оценки по месту.
// Это публичный API модуля, используемый модулем карты для отображения
// сводных характеристик места в карточке.
type Aggregate struct {
	PlaceID    uint64  `json:"place_id"`
	AvgNoise   float64 `json:"avg_noise"`
	AvgLight   float64 `json:"avg_light"`
	AvgCrowd   float64 `json:"avg_crowd"`
	AvgSmell   float64 `json:"avg_smell"`
	AvgVisual  float64 `json:"avg_visual"`
	OverallAvg float64 `json:"overall_avg"`
	ReviewsCnt int     `json:"reviews_count"`
}

func (s *ReviewService) Aggregate(placeID uint64) (Aggregate, error) {
	var a Aggregate
	a.PlaceID = placeID
	err := s.db.Model(&models.Review{}).
		Select(`COALESCE(AVG(noise),0) AS avg_noise,
			COALESCE(AVG(light),0) AS avg_light,
			COALESCE(AVG(crowd),0) AS avg_crowd,
			COALESCE(AVG(smell),0) AS avg_smell,
			COALESCE(AVG(visual),0) AS avg_visual,
			COALESCE(AVG((noise+light+crowd+smell+visual)/5.0),0) AS overall_avg,
			COUNT(*) AS reviews_cnt`).
		Where("place_id = ?", placeID).
		Scan(&a).Error
	return a, err
}

func (s *ReviewService) AddFavorite(userID, placeID uint64) error {
	var place models.PlaceRef
	if err := s.db.First(&place, placeID).Error; err != nil {
		return ErrPlaceMissing
	}
	fav := models.Favorite{UserID: userID, PlaceID: placeID}
	return s.db.Where("user_id = ? AND place_id = ?", userID, placeID).
		FirstOrCreate(&fav).Error
}

func (s *ReviewService) RemoveFavorite(userID, placeID uint64) error {
	return s.db.Where("user_id = ? AND place_id = ?", userID, placeID).
		Delete(&models.Favorite{}).Error
}

func (s *ReviewService) ListFavorites(userID uint64) ([]models.PlaceRef, error) {
	var places []models.PlaceRef
	err := s.db.
		Joins("JOIN favorites f ON f.place_id = places.id").
		Where("f.user_id = ?", userID).
		Order("f.created_at DESC").
		Find(&places).Error
	return places, err
}

func (s *ReviewService) IsFavorite(userID, placeID uint64) (bool, error) {
	var n int64
	err := s.db.Model(&models.Favorite{}).
		Where("user_id = ? AND place_id = ?", userID, placeID).
		Count(&n).Error
	return n > 0, err
}