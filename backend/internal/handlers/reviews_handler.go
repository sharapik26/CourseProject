package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/middleware"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/models"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/services"
)

type ReviewsHandler struct {
	reviews *services.ReviewService
}

func NewReviewsHandler(s *services.ReviewService) *ReviewsHandler {
	return &ReviewsHandler{reviews: s}
}

type reviewReq struct {
	Text   string `json:"text"`
	Noise  int    `json:"noise"  binding:"required,min=1,max=5"`
	Light  int    `json:"light"  binding:"required,min=1,max=5"`
	Crowd  int    `json:"crowd"  binding:"required,min=1,max=5"`
	Smell  int    `json:"smell"  binding:"required,min=1,max=5"`
	Visual int    `json:"visual" binding:"required,min=1,max=5"`
}

// Create — POST /api/places/:id/reviews
func (h *ReviewsHandler) Create(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	placeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	var req reviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := &models.Review{
		PlaceID: placeID,
		UserID:  uid,
		Text:    req.Text,
		Noise:   req.Noise,
		Light:   req.Light,
		Crowd:   req.Crowd,
		Smell:   req.Smell,
		Visual:  req.Visual,
	}
	if err := h.reviews.Create(r); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, r)
}

// Update — PUT /api/reviews/:id
func (h *ReviewsHandler) Update(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	rid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	var req reviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := &models.Review{
		ID: rid, Text: req.Text, Noise: req.Noise, Light: req.Light,
		Crowd: req.Crowd, Smell: req.Smell, Visual: req.Visual,
	}
	if err := h.reviews.Update(uid, r); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

// Delete — DELETE /api/reviews/:id
func (h *ReviewsHandler) Delete(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	rid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	if err := h.reviews.Delete(uid, rid); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListByPlace — GET /api/places/:id/reviews
func (h *ReviewsHandler) ListByPlace(c *gin.Context) {
	placeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	rs, err := h.reviews.ListByPlace(placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rs, "count": len(rs)})
}

// ListMyReviews — GET /api/users/me/reviews
func (h *ReviewsHandler) ListMyReviews(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	rs, err := h.reviews.ListByUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rs, "count": len(rs)})
}

// AggregateForPlace — GET /api/places/:id/aggregate
// Публичный API для модуля карты: возвращает агрегированные сенсорные оценки.
func (h *ReviewsHandler) AggregateForPlace(c *gin.Context) {
	placeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	a, err := h.reviews.Aggregate(placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// AddFavorite — POST /api/favorites/:id
func (h *ReviewsHandler) AddFavorite(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	placeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	if err := h.reviews.AddFavorite(uid, placeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// RemoveFavorite — DELETE /api/favorites/:id
func (h *ReviewsHandler) RemoveFavorite(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	placeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	if err := h.reviews.RemoveFavorite(uid, placeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListFavorites — GET /api/users/me/favorites
func (h *ReviewsHandler) ListFavorites(c *gin.Context) {
	uid, _ := middleware.CurrentUserID(c)
	places, err := h.reviews.ListFavorites(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": places, "count": len(places)})
}