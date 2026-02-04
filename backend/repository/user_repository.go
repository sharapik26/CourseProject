package repository

import (
	"database/sql"
	"fmt"
	"time"

	"sensory-navigator/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash, username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, username)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, username, avatar_url, birth_date, created_at, updated_at
	`, email, passwordHash, username).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username,
		&user.AvatarURL, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, username, avatar_url, birth_date, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username,
		&user.AvatarURL, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, username, avatar_url, birth_date, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username,
		&user.AvatarURL, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(id int64, username *string, avatarURL *string, birthDate *string) (*models.User, error) {
	// Build dynamic update query
	query := "UPDATE users SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argNum := 1

	if username != nil {
		query += fmt.Sprintf(", username = $%d", argNum)
		args = append(args, *username)
		argNum++
	}
	if avatarURL != nil {
		query += fmt.Sprintf(", avatar_url = $%d", argNum)
		args = append(args, *avatarURL)
		argNum++
	}
	if birthDate != nil {
		query += fmt.Sprintf(", birth_date = $%d", argNum)
		args = append(args, *birthDate)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, id)
	query += " RETURNING id, email, password_hash, username, avatar_url, birth_date, created_at, updated_at"

	user := &models.User{}
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username,
		&user.AvatarURL, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdatePassword(id int64, passwordHash string) error {
	_, err := r.db.Exec(`
		UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2
	`, passwordHash, id)
	return err
}

func (r *UserRepository) CreatePasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)
	return err
}

func (r *UserRepository) FindPasswordResetToken(token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow(`
		SELECT user_id FROM password_reset_tokens 
		WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP AND used = FALSE
	`, token).Scan(&userID)
	return userID, err
}

func (r *UserRepository) MarkPasswordResetTokenUsed(token string) error {
	_, err := r.db.Exec(`UPDATE password_reset_tokens SET used = TRUE WHERE token = $1`, token)
	return err
}

func (r *UserRepository) SaveRefreshToken(userID int64, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)
	return err
}

func (r *UserRepository) FindRefreshToken(token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow(`
		SELECT user_id FROM refresh_tokens 
		WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP AND revoked = FALSE
	`, token).Scan(&userID)
	return userID, err
}

func (r *UserRepository) RevokeRefreshToken(token string) error {
	_, err := r.db.Exec(`UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`, token)
	return err
}

func (r *UserRepository) RevokeAllUserRefreshTokens(userID int64) error {
	_, err := r.db.Exec(`UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`, userID)
	return err
}

