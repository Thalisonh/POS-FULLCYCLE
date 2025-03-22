package middleware

import (
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/limiter"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/storage"
	"github.com/gin-gonic/gin"
)

func RateLimiter(configs *configs.Config, redis storage.StorageInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		isBlock, statusCode, err := limiter.Block(c, configs, redis)
		if isBlock {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			c.Abort()

			return
		}

		c.Next()
	}
}
