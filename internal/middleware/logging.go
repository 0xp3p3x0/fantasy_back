package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("requestID", requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		start := time.Now()
		c.Next()

		latency := time.Since(start)
		log.Printf("request_id=%s method=%s path=%s status=%d latency=%s", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}
