package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"sensory-navigator/config"
	"sensory-navigator/models"
	"sensory-navigator/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.JWTConfig
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repository.UserRepository, jwtConfig *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   jwtConfig,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	return s.userRepo.Create(req.Email, string(hashedPassword), req.Username)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.User, *TokenPair, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Generate tokens
	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) GenerateTokenPair(userID int64) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	expiresAt := time.Now().Add(s.config.RefreshExpiry)
	if err := s.userRepo.SaveRefreshToken(userID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, error) {
	// Find refresh token in database
	userID, err := s.userRepo.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Revoke old refresh token
	s.userRepo.RevokeRefreshToken(refreshToken)

	// Generate new token pair
	return s.GenerateTokenPair(userID)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if user exists
		return nil
	}

	// Generate reset token
	token, err := s.generateRefreshToken()
	if err != nil {
		return err
	}

	// Save token with 1 hour expiry
	expiresAt := time.Now().Add(time.Hour)
	if err := s.userRepo.CreatePasswordResetToken(user.ID, token, expiresAt); err != nil {
		return err
	}

	// TODO: Send email with reset link
	// For now, just log the token (in production, send email)
	
	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Find token
	userID, err := s.userRepo.FindPasswordResetToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return err
	}

	// Mark token as used
	s.userRepo.MarkPasswordResetTokenUsed(token)

	// Revoke all refresh tokens
	s.userRepo.RevokeAllUserRefreshTokens(userID)

	return nil
}

func (s *AuthService) generateAccessToken(userID int64) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

