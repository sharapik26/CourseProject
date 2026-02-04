package repository

import (
	"database/sql"

	"sensory-navigator/models"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, placeID int64) (*models.Favorite, error) {
	favorite := &models.Favorite{}
	err := r.db.QueryRow(`
		INSERT INTO favorites (user_id, place_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, place_id) DO UPDATE SET user_id = favorites.user_id
		RETURNING id, user_id, place_id, created_at
	`, userID, placeID).Scan(
		&favorite.ID, &favorite.UserID, &favorite.PlaceID, &favorite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *FavoriteRepository) Remove(userID, placeID int64) error {
	_, err := r.db.Exec(`DELETE FROM favorites WHERE user_id = $1 AND place_id = $2`, userID, placeID)
	return err
}

func (r *FavoriteRepository) FindByUserID(userID int64, limit, offset int) ([]*models.FavoriteResponse, error) {
	rows, err := r.db.Query(`
		SELECT f.id, f.place_id, f.created_at, p.name, p.address, p.category
		FROM favorites f
		JOIN places p ON f.place_id = p.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []*models.FavoriteResponse
	for rows.Next() {
		fav := &models.FavoriteResponse{}
		var address, category sql.NullString
		err := rows.Scan(&fav.ID, &fav.PlaceID, &fav.CreatedAt, &fav.PlaceName, &address, &category)
		if err != nil {
			return nil, err
		}
		if address.Valid {
			fav.Address = address.String
		}
		if category.Valid {
			fav.Category = category.String
		}
		favorites = append(favorites, fav)
	}
	return favorites, nil
}

func (r *FavoriteRepository) Exists(userID, placeID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND place_id = $2)
	`, userID, placeID).Scan(&exists)
	return exists, err
}

func (r *FavoriteRepository) CountByUserID(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM favorites WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

