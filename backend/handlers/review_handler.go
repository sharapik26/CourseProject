package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sensory-navigator/models"
	"sensory-navigator/repository"
)

type ReviewHandler struct {
	reviewRepo *repository.ReviewRepository
}

func NewReviewHandler(reviewRepo *repository.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo}
}

// GET /api/places/:id/reviews
func (h *ReviewHandler) GetPlaceReviews(c *gin.Context) {
	placeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, err := h.reviewRepo.FindByPlaceID(placeID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"limit":   limit,
		"offset":  offset,
	})
}

// POST /api/places/:id/reviews
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	placeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place ID"})
		return
	}

	// Check if user already reviewed this place
	exists, _ := h.reviewRepo.ExistsByUserAndPlace(userID, placeID)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "you have already reviewed this place"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewRepo.Create(userID, placeID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review.ToResponse())
}

// PUT /api/reviews/:id
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	// Check ownership
	existingReview, err := h.reviewRepo.FindByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if existingReview.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only edit your own reviews"})
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewRepo.Update(reviewID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review.ToResponse())
}

// DELETE /api/reviews/:id
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	// Check ownership
	existingReview, err := h.reviewRepo.FindByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if existingReview.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own reviews"})
		return
	}

	if err := h.reviewRepo.Delete(reviewID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}

