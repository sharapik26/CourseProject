package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sensory-navigator/models"
	"sensory-navigator/repository"
)

type UserHandler struct {
	userRepo     *repository.UserRepository
	reviewRepo   *repository.ReviewRepository
	favoriteRepo *repository.FavoriteRepository
}

func NewUserHandler(userRepo *repository.UserRepository, reviewRepo *repository.ReviewRepository, favoriteRepo *repository.FavoriteRepository) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		reviewRepo:   reviewRepo,
		favoriteRepo: favoriteRepo,
	}
}

// GET /api/users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID")

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// PUT /api/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.Update(userID, req.Username, req.AvatarURL, req.BirthDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// GET /api/users/me/reviews
func (h *UserHandler) GetMyReviews(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, err := h.reviewRepo.FindByUserID(userID, limit, offset)
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

// GET /api/users/me/favorites
func (h *UserHandler) GetMyFavorites(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	favorites, err := h.favoriteRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch favorites"})
		return
	}

	count, _ := h.favoriteRepo.CountByUserID(userID)

	c.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
		"total":     count,
		"limit":     limit,
		"offset":    offset,
	})
}

