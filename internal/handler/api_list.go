package handler

import (
	"back/internal/model"
	"back/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIListHandler struct {
	svc *service.APIListService
}

func NewAPIListHandler(svc *service.APIListService) *APIListHandler {
	return &APIListHandler{svc: svc}
}

func (h *APIListHandler) Create(c *gin.Context) {
	var req service.CreateAPIListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	row, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, model.APIResponse{Success: true, Data: row})
}

func (h *APIListHandler) List(c *gin.Context) {
	rows, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Success: true, Data: rows})
}

func (h *APIListHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: "invalid id"})
		return
	}
	row, err := h.svc.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.APIResponse{Success: false, Message: "api list not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Success: true, Data: row})
}

func (h *APIListHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: "invalid id"})
		return
	}
	var req service.UpdateAPIListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	row, err := h.svc.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.APIResponse{Success: false, Message: "api list not found"})
			return
		}
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Success: true, Data: row})
}

func (h *APIListHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: "invalid id"})
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.APIResponse{Success: false, Message: "api list not found"})
			return
		}
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Success: true, Message: "deleted"})
}
