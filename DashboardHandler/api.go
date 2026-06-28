package main

import (
	"net/http"

	"github.com/SergeStormTeam/dashboard-handler/dashboard"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ViewLiveData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	dashboard.AddWebSocketConnection(conn)
	defer func() {
		dashboard.RemoveWebSocketConnection(conn)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
