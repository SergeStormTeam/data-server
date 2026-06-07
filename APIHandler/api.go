package main

import (
	"encoding/json"
	"net/http"

	"github.com/WeatherGod3218/serge-api-handler/logging"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type ProbeData struct {
	Timestamp     float64  `json:"timestamp"`
	CO2           *float64 `json:"co2"`
	Humidity      *float64 `json:"humidity"`
	Precipitation *float64 `json:"precipitation"`
	Pressure      *float64 `json:"pressure"`
	VOC           *float64 `json:"voc"`
	WindSpeed     *float64 `json:"wind_speed"`
}

type DatabaseEntry struct {
	Timestamp     float64  `json:"timestamp"`
	SessionId     string   `json:"session_id"`
	Sequence      int      `json:"sequence"`
	CO2           *float64 `json:"co2"`
	Humidity      *float64 `json:"humidity"`
	Precipitation *float64 `json:"precipitation"`
	Pressure      *float64 `json:"pressure"`
	VOC           *float64 `json:"voc"`
	WindSpeed     *float64 `json:"wind_speed"`
}

type DatabaseResponse struct {
	SessionId string `json:"session_id"`
	Sequence  int    `json:"sequence"`
}

type DatabaseBackup struct {
	Data []DatabaseEntry `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpdateLiveData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType != 1 {
			return
		}

		var data ProbeData

		err = json.Unmarshal(message, &data)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "UpdateLiveData"}).Warn("Unable to unmarshel data from websocket!")
			break
		}

		// logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "UpdateLiveData"}).Info(fmt.Sprintf("Succesfully recieved new data! , Timestamp: %f, CO2: %v, Humidity: %v, Precipitation: %v Pressure: %v VOC: %v WindSpeed: %v",
		// 	data.Timestamp,
		// 	data.CO2,
		// 	data.Humidity,
		// 	data.Precipitation,
		// 	data.Pressure,
		// 	data.VOC,
		// 	data.WindSpeed,
		// ))
	}
}

func UpdateDatabase(c *gin.Context) {

	var req DatabaseBackup

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var rows []DatabaseResponse

	rows = make([]DatabaseResponse, 3)
	c.JSON(http.StatusOK, gin.H{
		"updated": rows,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
