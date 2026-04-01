package handler

import (
	"back/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentService *service.AgentService
	authService  *service.AuthService
}

func NewAgentHandler(agentService *service.AgentService, authService *service.AuthService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		authService:  authService,
	}
}

func (h *AgentHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user.Role != "agent" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only agents can use this endpoint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// CreateAgent godoc
// @Summary Create agent (admin only)
// @Tags agent-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateAgentRequest true "Create agent payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/admin/agents [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req service.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.agentService.CreateAgent(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"agent": agent})
}

// ListAgents godoc
// @Summary List all agents (admin only)
// @Tags agent-admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/admin/agents [get]
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents, err := h.agentService.ListAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// UpdateAgent godoc
// @Summary Update agent information (admin only)
// @Tags agent-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Agent ID"
// @Param body body service.UpdateAgentRequest true "Update agent payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/admin/agents/{id} [put]
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent id"})
		return
	}

	var req service.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.agentService.UpdateAgent(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

// DeleteAgent godoc
// @Summary Delete an agent (admin only)
// @Tags agent-admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Agent ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/admin/agents/{id} [delete]
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent id"})
		return
	}

	if err := h.agentService.DeleteAgent(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "agent deleted"})
}
