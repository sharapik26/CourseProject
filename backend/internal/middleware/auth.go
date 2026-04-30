package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/auth"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
)

const UserContextKey = "current_user_id"

// RequireAuth требует наличия валидного JWT-токена в заголовке Authorization.
// Этот middleware экспортируется как часть публичного API модуля и используется
// другими модулями (модуль карты Атаханова Н. Р.) для защиты эндпоинтов записи.
func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректный формат токена"})
			return
		}
		claims, err := auth.ParseToken(parts[1], cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(UserContextKey, claims.UserID)
		c.Next()
	}
}

// CurrentUserID возвращает идентификатор пользователя из контекста запроса.
func CurrentUserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(UserContextKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}