package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/database"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/routes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Интеграционные тесты выполняются только при наличии переменной TEST_DB_DSN,
// указывающей на пустую тестовую базу PostgreSQL.
func openTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN не задан — пропускаем интеграционный тест")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	return db
}

func TestRegisterAndLogin(t *testing.T) {
	db := openTestDB(t)
	cfg := &config.Config{
		JWTSecret: "test-secret", JWTAccessTTL: time.Hour, BCryptCost: 4,
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	r := routes.NewRouter(db, cfg)

	body := `{"email":"alice@example.com","username":"alice","password":"secret123","display_name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NotEmpty(t, registerResp.Token)

	loginBody := `{"email":"alice@example.com","password":"secret123"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}