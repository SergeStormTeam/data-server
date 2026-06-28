package main

import (
	"github.com/SergeStormTeam/dashboard-handler/logging"
	"github.com/SergeStormTeam/dashboard-handler/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	err := redis.InitRedis()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitRedis"}).Fatal("Failed to init redis!")
	}
	go redis.SubscribeToZephyrApi()

	router := gin.Default()

	group := router.Group("/")

	group.Use(redis.RedisRateLimiter(2, 50))

	group.GET("/health", HealthCheck)

	// Websocket!
	group.GET("/refresh-live", ViewLiveData)

	router.Run(":5632")
}
