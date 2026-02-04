package repository

import (
	"database/sql"

	"sensory-navigator/models"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(userID, placeID int64, req *models.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.QueryRow(`
		INSERT INTO reviews (user_id, place_id, text, sensory_rating, lighting_rating, 
			sound_level_rating, crowding_rating, accessibility_rating, overall_rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, place_id, text, sensory_rating, lighting_rating, 
			sound_level_rating, crowding_rating, accessibility_rating, overall_rating, created_at, updated_at
	`, userID, placeID, req.Text, req.SensoryRating, req.LightingRating,
		req.SoundLevelRating, req.CrowdingRating, req.AccessibilityRating, req.OverallRating).Scan(
		&review.ID, &review.UserID, &review.PlaceID, &review.Text,
		&review.SensoryRating, &review.LightingRating, &review.SoundLevelRating,
		&review.CrowdingRating, &review.AccessibilityRating, &review.OverallRating,
		&review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) FindByID(id int64) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.QueryRow(`
		SELECT id, user_id, place_id, text, sensory_rating, lighting_rating,
			sound_level_rating, crowding_rating, accessibility_rating, overall_rating, created_at, updated_at
		FROM reviews WHERE id = $1
	`, id).Scan(
		&review.ID, &review.UserID, &review.PlaceID, &review.Text,
		&review.SensoryRating, &review.LightingRating, &review.SoundLevelRating,
		&review.CrowdingRating, &review.AccessibilityRating, &review.OverallRating,
		&review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) FindByPlaceID(placeID int64, limit, offset int) ([]*models.ReviewResponse, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.user_id, r.place_id, r.text, r.sensory_rating, r.lighting_rating,
			r.sound_level_rating, r.crowding_rating, r.accessibility_rating, r.overall_rating, 
			r.created_at, r.updated_at, u.username
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.place_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`, placeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.ReviewResponse
	for rows.Next() {
		review := &models.Review{}
		var username string
		err := rows.Scan(
			&review.ID, &review.UserID, &review.PlaceID, &review.Text,
			&review.SensoryRating, &review.LightingRating, &review.SoundLevelRating,
			&review.CrowdingRating, &review.AccessibilityRating, &review.OverallRating,
			&review.CreatedAt, &review.UpdatedAt, &username,
		)
		if err != nil {
			return nil, err
		}
		resp := review.ToResponse()
		resp.Username = username
		reviews = append(reviews, &resp)
	}
	return reviews, nil
}

func (r *ReviewRepository) FindByUserID(userID int64, limit, offset int) ([]*models.ReviewResponse, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.user_id, r.place_id, r.text, r.sensory_rating, r.lighting_rating,
			r.sound_level_rating, r.crowding_rating, r.accessibility_rating, r.overall_rating, 
			r.created_at, r.updated_at, p.name
		FROM reviews r
		JOIN places p ON r.place_id = p.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.ReviewResponse
	for rows.Next() {
		review := &models.Review{}
		var placeName string
		err := rows.Scan(
			&review.ID, &review.UserID, &review.PlaceID, &review.Text,
			&review.SensoryRating, &review.LightingRating, &review.SoundLevelRating,
			&review.CrowdingRating, &review.AccessibilityRating, &review.OverallRating,
			&review.CreatedAt, &review.UpdatedAt, &placeName,
		)
		if err != nil {
			return nil, err
		}
		resp := review.ToResponse()
		resp.PlaceName = placeName
		reviews = append(reviews, &resp)
	}
	return reviews, nil
}

func (r *ReviewRepository) Update(id int64, req *models.UpdateReviewRequest) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.QueryRow(`
		UPDATE reviews SET 
			text = COALESCE($1, text),
			sensory_rating = COALESCE($2, sensory_rating),
			lighting_rating = COALESCE($3, lighting_rating),
			sound_level_rating = COALESCE($4, sound_level_rating),
			crowding_rating = COALESCE($5, crowding_rating),
			accessibility_rating = COALESCE($6, accessibility_rating),
			overall_rating = COALESCE($7, overall_rating),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING id, user_id, place_id, text, sensory_rating, lighting_rating,
			sound_level_rating, crowding_rating, accessibility_rating, overall_rating, created_at, updated_at
	`, req.Text, req.SensoryRating, req.LightingRating, req.SoundLevelRating,
		req.CrowdingRating, req.AccessibilityRating, req.OverallRating, id).Scan(
		&review.ID, &review.UserID, &review.PlaceID, &review.Text,
		&review.SensoryRating, &review.LightingRating, &review.SoundLevelRating,
		&review.CrowdingRating, &review.AccessibilityRating, &review.OverallRating,
		&review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM reviews WHERE id = $1`, id)
	return err
}

func (r *ReviewRepository) ExistsByUserAndPlace(userID, placeID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND place_id = $2)
	`, userID, placeID).Scan(&exists)
	return exists, err
}

