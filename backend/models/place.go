package models

import (
	"database/sql"
	"time"
)

type Place struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Address   sql.NullString `json:"address,omitempty"`
	Latitude  sql.NullFloat64 `json:"latitude,omitempty"`
	Longitude sql.NullFloat64 `json:"longitude,omitempty"`
	Category  sql.NullString `json:"category,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type PlaceResponse struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Category  *string  `json:"category,omitempty"`
	CreatedAt string   `json:"created_at"`
}

func (p *Place) ToResponse() PlaceResponse {
	resp := PlaceResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	}

	if p.Address.Valid {
		resp.Address = &p.Address.String
	}
	if p.Latitude.Valid {
		resp.Latitude = &p.Latitude.Float64
	}
	if p.Longitude.Valid {
		resp.Longitude = &p.Longitude.Float64
	}
	if p.Category.Valid {
		resp.Category = &p.Category.String
	}

	return resp
}

