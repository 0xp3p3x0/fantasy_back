package handler

import (
	"back/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

type updateProfileRequest struct {
	CallbackURL string `json:"callback_url" binding:"required,url"`
}

type regenerateTokenRequest struct {
	Username string `json:"username" binding:"required"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfileById(c *gin.Context) {
	profile, err := h.profileService.GetProfileById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetProfileByCode(c *gin.Context) {
	profile, err := h.profileService.GetProfileByCode(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.profileService.UpdateCallbackURLByID(userID, req.CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateCallbackURL(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.profileService.UpdateCallbackURLByID(userID, req.CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.ChangePasswordByID(userID, req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func getUserIDFromClaims(c *gin.Context) (uint, error) {
	claimsAny, ok := c.Get("claims")
	if !ok {
		return 0, errors.New("missing auth claims")
	}

	claims, ok := claimsAny.(*jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid auth claims")
	}

	sub, ok := (*claims)["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid user id in token")
	}

	return uint(sub), nil
}

func (h *ProfileHandler) RegenerateToken(c *gin.Context) {
	var req regenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := h.profileService.RegenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}
