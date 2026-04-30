package auth_test

import (
	"testing"
	"time"

	"github.com/GyBJluHv2/sensory-navigator-users/internal/auth"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := auth.GeneratePasswordHash("super-secret-123", 4)
	assert.NoError(t, err)
	assert.True(t, auth.CheckPassword("super-secret-123", hash))
	assert.False(t, auth.CheckPassword("wrong", hash))
}

func TestJWTRoundtrip(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTAccessTTL: time.Hour}
	token, exp, err := auth.IssueToken(42, cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, exp.After(time.Now()))

	claims, err := auth.ParseToken(token, cfg)
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), claims.UserID)
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	cfgIssue := &config.Config{JWTSecret: "secret-A", JWTAccessTTL: time.Hour}
	cfgVerify := &config.Config{JWTSecret: "secret-B", JWTAccessTTL: time.Hour}

	token, _, err := auth.IssueToken(1, cfgIssue)
	assert.NoError(t, err)

	_, err = auth.ParseToken(token, cfgVerify)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}