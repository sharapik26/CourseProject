package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/auth"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/config"
	"github.com/GyBJluHv2/sensory-navigator-users/internal/services"
)

type AuthHandler struct {
	users         *services.UserService
	verifications *services.VerificationService
	cfg           *config.Config
}

func NewAuthHandler(users *services.UserService, verifications *services.VerificationService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{users: users, verifications: verifications, cfg: cfg}
}

type registerReq struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=64"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tokenResp struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	User      any    `json:"user"`
}

// Register — POST /api/auth/register (одношаговый flow для seed/тестов).
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.users.Register(services.RegisterInput{
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	token, exp, err := auth.IssueToken(user.ID, h.cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tokenResp{
		Token:     token,
		ExpiresAt: exp.Format("2006-01-02T15:04:05Z07:00"),
		User:      user,
	})
}

// Login — POST /api/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.users.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, exp, err := auth.IssueToken(user.ID, h.cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokenResp{
		Token:     token,
		ExpiresAt: exp.Format("2006-01-02T15:04:05Z07:00"),
		User:      user,
	})
}

// RequestRegister — POST /api/auth/register-request.
//
// Шаг 1 регистрации: проверка формата email, отправка 6-значного кода.
func (h *AuthHandler) RequestRegister(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exp, err := h.verifications.Request(services.RequestRegisterInput{
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "code_sent",
		"expires_at": exp.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type confirmReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code"  binding:"required,len=6"`
}

// ConfirmRegister — POST /api/auth/register-confirm.
func (h *AuthHandler) ConfirmRegister(c *gin.Context) {
	var req confirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.verifications.Confirm(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, exp, err := auth.IssueToken(user.ID, h.cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tokenResp{
		Token:     token,
		ExpiresAt: exp.Format("2006-01-02T15:04:05Z07:00"),
		User:      user,
	})
}

type resendReq struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendCode — POST /api/auth/resend-code.
func (h *AuthHandler) ResendCode(c *gin.Context) {
	var req resendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exp, err := h.verifications.Resend(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "code_sent",
		"expires_at": exp.Format("2006-01-02T15:04:05Z07:00"),
	})
}