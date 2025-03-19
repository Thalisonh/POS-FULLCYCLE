package server

import (
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/middleware"
	"github.com/gin-gonic/gin"
)

func StartServer(configs *configs.Config) {
	r := gin.Default()

	// Injetar o middleware do rate limiter
	r.Use(middleware.RateLimiter(configs))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API funcionando"})
	})

	r.Run(":8080")
}
