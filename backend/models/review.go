package models

import (
	"database/sql"
	"time"
)

type Review struct {
	ID                  int64          `json:"id"`
	UserID              int64          `json:"user_id"`
	PlaceID             int64          `json:"place_id"`
	Text                sql.NullString `json:"text,omitempty"`
	SensoryRating       sql.NullInt32  `json:"sensory_rating,omitempty"`
	LightingRating      sql.NullInt32  `json:"lighting_rating,omitempty"`
	SoundLevelRating    sql.NullInt32  `json:"sound_level_rating,omitempty"`
	CrowdingRating      sql.NullInt32  `json:"crowding_rating,omitempty"`
	AccessibilityRating sql.NullInt32  `json:"accessibility_rating,omitempty"`
	OverallRating       sql.NullFloat64 `json:"overall_rating,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

type ReviewResponse struct {
	ID                  int64    `json:"id"`
	UserID              int64    `json:"user_id"`
	PlaceID             int64    `json:"place_id"`
	Text                *string  `json:"text,omitempty"`
	SensoryRating       *int     `json:"sensory_rating,omitempty"`
	LightingRating      *int     `json:"lighting_rating,omitempty"`
	SoundLevelRating    *int     `json:"sound_level_rating,omitempty"`
	CrowdingRating      *int     `json:"crowding_rating,omitempty"`
	AccessibilityRating *int     `json:"accessibility_rating,omitempty"`
	OverallRating       *float64 `json:"overall_rating,omitempty"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Username            string   `json:"username,omitempty"`
	PlaceName           string   `json:"place_name,omitempty"`
}

type CreateReviewRequest struct {
	Text                *string  `json:"text,omitempty"`
	SensoryRating       *int     `json:"sensory_rating,omitempty" binding:"omitempty,min=1,max=5"`
	LightingRating      *int     `json:"lighting_rating,omitempty" binding:"omitempty,min=1,max=5"`
	SoundLevelRating    *int     `json:"sound_level_rating,omitempty" binding:"omitempty,min=1,max=5"`
	CrowdingRating      *int     `json:"crowding_rating,omitempty" binding:"omitempty,min=1,max=5"`
	AccessibilityRating *int     `json:"accessibility_rating,omitempty" binding:"omitempty,min=1,max=5"`
	OverallRating       *float64 `json:"overall_rating,omitempty" binding:"omitempty,min=1,max=5"`
}

type UpdateReviewRequest struct {
	Text                *string  `json:"text,omitempty"`
	SensoryRating       *int     `json:"sensory_rating,omitempty" binding:"omitempty,min=1,max=5"`
	LightingRating      *int     `json:"lighting_rating,omitempty" binding:"omitempty,min=1,max=5"`
	SoundLevelRating    *int     `json:"sound_level_rating,omitempty" binding:"omitempty,min=1,max=5"`
	CrowdingRating      *int     `json:"crowding_rating,omitempty" binding:"omitempty,min=1,max=5"`
	AccessibilityRating *int     `json:"accessibility_rating,omitempty" binding:"omitempty,min=1,max=5"`
	OverallRating       *float64 `json:"overall_rating,omitempty" binding:"omitempty,min=1,max=5"`
}

func (r *Review) ToResponse() ReviewResponse {
	resp := ReviewResponse{
		ID:        r.ID,
		UserID:    r.UserID,
		PlaceID:   r.PlaceID,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}

	if r.Text.Valid {
		resp.Text = &r.Text.String
	}
	if r.SensoryRating.Valid {
		val := int(r.SensoryRating.Int32)
		resp.SensoryRating = &val
	}
	if r.LightingRating.Valid {
		val := int(r.LightingRating.Int32)
		resp.LightingRating = &val
	}
	if r.SoundLevelRating.Valid {
		val := int(r.SoundLevelRating.Int32)
		resp.SoundLevelRating = &val
	}
	if r.CrowdingRating.Valid {
		val := int(r.CrowdingRating.Int32)
		resp.CrowdingRating = &val
	}
	if r.AccessibilityRating.Valid {
		val := int(r.AccessibilityRating.Int32)
		resp.AccessibilityRating = &val
	}
	if r.OverallRating.Valid {
		resp.OverallRating = &r.OverallRating.Float64
	}

	return resp
}

