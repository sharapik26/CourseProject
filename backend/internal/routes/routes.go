package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/email"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/handlers"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/middleware"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/services"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	userSvc := services.NewUserService(db, cfg)
	reviewSvc := services.NewReviewService(db)

	// Email-отправитель: stub в логи или реальный SMTP по env-конфигу.
	var mailer email.Sender
	if cfg.SMTPHost == "" {
		mailer = email.LogSender{}
	} else {
		mailer = email.SMTPSender{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			User:     cfg.SMTPUser,
			Password: cfg.SMTPPassword,
			From:     cfg.SMTPFrom,
			UseTLS:   cfg.SMTPUseTLS,
		}
	}
	verificationSvc := services.NewVerificationService(db, cfg, userSvc, mailer)

	authH := handlers.NewAuthHandler(userSvc, verificationSvc, cfg)
	usersH := handlers.NewUsersHandler(userSvc)
	reviewsH := handlers.NewReviewsHandler(reviewSvc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")

	// Открытые маршруты авторизации
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/register-request", authH.RequestRegister)
	api.POST("/auth/register-confirm", authH.ConfirmRegister)
	api.POST("/auth/resend-code", authH.ResendCode)
	api.POST("/auth/login", authH.Login)

	// Чтение отзывов и агрегаты — публичные (используются модулем карты)
	api.GET("/places/:id/reviews", reviewsH.ListByPlace)
	api.GET("/places/:id/aggregate", reviewsH.AggregateForPlace)

	// Защищённые маршруты
	auth := api.Group("/")
	auth.Use(middleware.RequireAuth(cfg))
	{
		auth.GET("/users/me", usersH.Me)
		auth.PUT("/users/me", usersH.UpdateMe)
		auth.PUT("/users/me/password", usersH.ChangePassword)
		auth.GET("/users/me/reviews", reviewsH.ListMyReviews)
		auth.GET("/users/me/favorites", reviewsH.ListFavorites)

		auth.POST("/places/:id/reviews", reviewsH.Create)
		auth.PUT("/reviews/:id", reviewsH.Update)
		auth.DELETE("/reviews/:id", reviewsH.Delete)

		auth.POST("/favorites/:id", reviewsH.AddFavorite)
		auth.DELETE("/favorites/:id", reviewsH.RemoveFavorite)
	}

	return r
}