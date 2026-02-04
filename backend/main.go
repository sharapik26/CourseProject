package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"sensory-navigator/config"
	"sensory-navigator/database"
	"sensory-navigator/handlers"
	"sensory-navigator/middleware"
	"sensory-navigator/repository"
	"sensory-navigator/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.DB); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.GetDB())
	reviewRepo := repository.NewReviewRepository(database.GetDB())
	favoriteRepo := repository.NewFavoriteRepository(database.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo, &cfg.JWT)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, reviewRepo, favoriteRepo)
	reviewHandler := handlers.NewReviewHandler(reviewRepo)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteRepo)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.GET("/me/reviews", userHandler.GetMyReviews)
				users.GET("/me/favorites", userHandler.GetMyFavorites)
			}

			// Review routes
			protected.GET("/places/:id/reviews", reviewHandler.GetPlaceReviews)
			protected.POST("/places/:id/reviews", reviewHandler.CreateReview)
			protected.PUT("/reviews/:id", reviewHandler.UpdateReview)
			protected.DELETE("/reviews/:id", reviewHandler.DeleteReview)

			// Favorite routes
			favorites := protected.Group("/favorites")
			{
				favorites.POST("/:placeId", favoriteHandler.AddFavorite)
				favorites.DELETE("/:placeId", favoriteHandler.RemoveFavorite)
				favorites.GET("/:placeId/check", favoriteHandler.CheckFavorite)
			}
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

