package middleware

import (
	"net/http"

	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/storage"
	"github.com/gin-gonic/gin"
)

var redisStore = storage.NewRedisStorage("localhost:6379")

func RateLimiter(configs *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		key := "rate_limit:" + ip
		limit := configs.RateLimitIP

		if token != "" {
			key = "rate_limit:token:" + token
			limit = configs.RateLimitToken
		}

		count, err := redisStore.Increment(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "you have reached the maximum number of requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
