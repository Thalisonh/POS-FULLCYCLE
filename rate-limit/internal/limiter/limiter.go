package limiter

import (
	"errors"
	"net/http"
	"time"

	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/storage"
	"github.com/gin-gonic/gin"
)

var message = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func Block(c *gin.Context, configs *configs.Config, redis storage.StorageInterface) (bool, int, error) {
	ip := c.ClientIP()
	token := c.GetHeader("API_KEY")

	key := "rate_limit:" + ip
	blockKey := "block:" + ip
	limit := configs.RateLimitIP

	if token != "" {
		key = "rate_limit:token:" + token
		blockKey = "block:" + token
		limit = configs.RateLimitToken
	}

	isBlocked, err := redis.Get(c, blockKey)
	if err != nil {
		return true, http.StatusInternalServerError, err
	}

	if isBlocked > 0 {
		return true, http.StatusTooManyRequests, errors.New(message)
	}

	count, err := redis.Increment(c, key)
	if err != nil {
		return true, http.StatusInternalServerError, err
	}

	if count > limit {
		// set if not exist, will block only first request
		_, err := redis.SetNX(c, blockKey, 1, time.Duration(configs.BlockTime)*time.Second)
		if err != nil {
			return true, http.StatusInternalServerError, err
		}

		return true, http.StatusTooManyRequests, errors.New(message)
	}

	return false, http.StatusOK, nil
}
