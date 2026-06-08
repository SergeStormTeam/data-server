package redis

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/WeatherGod3218/serge-api-handler/logging"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

var client *redis.Client

func InitRedis() error {

	newClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})

	pong, err := newClient.Ping(ctx).Result()

	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err := newClient.Ping(context.Background()).Err()
		if err == nil {
			client = newClient
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	logging.Logger.WithFields(logrus.Fields{"module": "redis", "method": "InitRedis"}).Info(fmt.Sprintf("Pinged Redis!: %s", pong))

	return nil
}

func RedisRateLimiter(rate float64, capacity float64) gin.HandlerFunc {

	limiter := NewTokenBucket(rate, capacity)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		allowed, tokens, err := limiter.Allow(c, ip)

		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "RedisRateLimiter"}).Warn(fmt.Sprintf("Failure in the redis cache %v", err))
		} else if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests!",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%v", tokens))

		c.Next()
	}
}
