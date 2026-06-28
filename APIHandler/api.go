package main

import (
	"net/http"

	"github.com/SergeStormTeam/api-handler/database"
	"github.com/SergeStormTeam/api-handler/logging"
	"github.com/SergeStormTeam/api-handler/redis"
	"github.com/SergeStormTeam/api-handler/types"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type ProbeData struct {
	Timestamp     float64  `json:"timestamp"`
	Temperature   *float64 `json:"temperature"`
	CO2           *float64 `json:"co2"`
	Humidity      *float64 `json:"humidity"`
	Precipitation *float64 `json:"precipitation"`
	Pressure      *float64 `json:"pressure"`
	VOC           *float64 `json:"voc"`
	WindSpeed     *float64 `json:"wind_speed"`
}

func LiveData(c *gin.Context) {
	var payload types.ZephyrUpdate
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to publish data"})
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "LiveData"}).Warn("Unable to parse data")
		return
	}

	err = redis.PublishToDashboard(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to publish data"})
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "LiveData"}).Warn("Unable to publish information to websocket!")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func UpdateDatabase(c *gin.Context) {

	var req database.DatabaseBackup

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to update in database"})
		return
	}

	data_rows, err := database.AddDataToDatabase(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update in database"})
		return
	}

	event_rows, err := database.AddEventsToDatabase(req.Events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   data_rows,
		"events": event_rows,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
