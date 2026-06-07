package main

import (
	"os"

	"github.com/WeatherGod3218/serge-api-handler/authorization"
	"github.com/WeatherGod3218/serge-api-handler/database"
	"github.com/WeatherGod3218/serge-api-handler/logging"
	"github.com/WeatherGod3218/serge-api-handler/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	err := redis.InitRedis()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitRedis"}).Fatal("Failed to init redis!")
	}

	err = database.InitDB()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitRedis"}).Fatal("Failed to init database!")
	}

	router := gin.Default()

	group := router.Group("/")

	group.Use(redis.RedisRateLimiter(2, 50))
	group.Use(authorization.CheckAuthorization(os.Getenv("SERVER_API_KEY")))

	group.GET("/health", HealthCheck)
	group.POST("/backup-data", UpdateDatabase)

	// Websocket!
	group.GET("/refresh-live", UpdateLiveData)

	router.Run(":8080")
}
