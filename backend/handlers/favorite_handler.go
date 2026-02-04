package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sensory-navigator/repository"
)

type FavoriteHandler struct {
	favoriteRepo *repository.FavoriteRepository
}

func NewFavoriteHandler(favoriteRepo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{favoriteRepo: favoriteRepo}
}

// POST /api/favorites/:placeId
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	placeID, err := strconv.ParseInt(c.Param("placeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place ID"})
		return
	}

	favorite, err := h.favoriteRepo.Add(userID, placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favorite"})
		return
	}

	c.JSON(http.StatusCreated, favorite.ToResponse())
}

// DELETE /api/favorites/:placeId
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	placeID, err := strconv.ParseInt(c.Param("placeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place ID"})
		return
	}

	if err := h.favoriteRepo.Remove(userID, placeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}

// GET /api/favorites/:placeId/check
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	userID := c.GetInt64("userID")
	
	placeID, err := strconv.ParseInt(c.Param("placeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place ID"})
		return
	}

	exists, _ := h.favoriteRepo.Exists(userID, placeID)

	c.JSON(http.StatusOK, gin.H{"is_favorite": exists})
}

