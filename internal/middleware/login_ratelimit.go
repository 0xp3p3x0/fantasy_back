package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// LoginRateLimit limits POST /auth/login per client IP using Redis. No-op when rdb is nil.
func LoginRateLimit(rdb *redis.Client, maxPerWindow int, window time.Duration) gin.HandlerFunc {
	if rdb == nil {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key := "ratelimit:login:" + c.ClientIP()

		n, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if n == 1 {
			_ = rdb.Expire(ctx, key, window).Err()
		}
		if n > int64(maxPerWindow) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts, try again later"})
			return
		}
		c.Next()
	}
}
