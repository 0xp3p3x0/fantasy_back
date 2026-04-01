package middleware

import (
	"errors"
	"log"
	"net/http"

	"back/internal/model"
	"back/internal/service"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		lastErr := c.Errors.Last().Err
		statusCode := http.StatusInternalServerError
		message := "internal server error"

		var appErr *service.AppError
		if errors.As(lastErr, &appErr) {
			statusCode = appErr.StatusCode
			message = appErr.Message
		} else {
			message = lastErr.Error()
		}

		requestID, _ := c.Get("requestID")
		log.Printf("request_id=%v error=%v", requestID, lastErr)

		c.AbortWithStatusJSON(statusCode, model.APIResponse{
			Success: false,
			Message: message,
		})
	}
}
