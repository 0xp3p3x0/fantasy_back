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

	res, err := h.casinoService.GetGameURL(c.Request.Context(), &req)
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

// Callback is the single webhook URL for all provider wallet calls.
// Body must include "method": "GetBalance" | "ChangeBalance" | "UpdateDetail" (case-insensitive);
// other fields depend on the method. Response is always { status, msg, balance? } per vendor spec.
func (h *CasinoHandler) Callback(c *gin.Context) {
	var req model.CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.ProviderCallbackResponse{
			Status: model.ProvCallbackBadRequest,
			Msg:    "INVALID_JSON",
		})
		return
	}
	out := h.casinoService.HandleProviderCallback(c.Request.Context(), &req)
	c.JSON(http.StatusOK, out)
}
