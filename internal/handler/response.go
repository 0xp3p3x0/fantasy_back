package handler

import (
	"back/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, model.APIResponse{
		Success: true,
		Data:    data,
	})
}

func fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, model.APIResponse{
		Success: false,
		Error: &model.APIError{
			Code:    code,
			Message: message,
		},
	})
}

func badRequest(c *gin.Context, message string) {
	fail(c, http.StatusBadRequest, "bad_request", message)
}

func unauthorized(c *gin.Context, message string) {
	fail(c, http.StatusUnauthorized, "unauthorized", message)
}

func forbidden(c *gin.Context, message string) {
	fail(c, http.StatusForbidden, "forbidden", message)
}
