package models

import "time"

type Favorite struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PlaceID   int64     `json:"place_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FavoriteResponse struct {
	ID        int64  `json:"id"`
	PlaceID   int64  `json:"place_id"`
	PlaceName string `json:"place_name,omitempty"`
	Address   string `json:"address,omitempty"`
	Category  string `json:"category,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (f *Favorite) ToResponse() FavoriteResponse {
	return FavoriteResponse{
		ID:        f.ID,
		PlaceID:   f.PlaceID,
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
	}
}

