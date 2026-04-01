package handler

import (
	"back/internal/model"
	"back/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CasinoHandler struct {
	casinoService *service.CasinoService
}

func NewCasinoHandler(casinoService *service.CasinoService) *CasinoHandler {
	return &CasinoHandler{casinoService: casinoService}
}

func (h *CasinoHandler) GetGameURL(c *gin.Context) {
	var req model.GetGameURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	res, err := h.casinoService.GetGameURL(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if appErr, ok := err.(*service.AppError); ok {
			statusCode = appErr.StatusCode
		}
		c.JSON(statusCode, model.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    res,
	})
}

