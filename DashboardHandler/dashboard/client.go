package dashboard

import (
	"sync"

	"github.com/SergeStormTeam/dashboard-handler/logging"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	conn *websocket.Conn
	send chan any
	once sync.Once
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.send)
		c.conn.Close()
	})
}

func (c *Client) SendLoop() {
	defer c.Close()

	for payload := range c.send {
		err := c.conn.WriteJSON(payload)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "dashboard", "method": "sendLoop"}).Warn("Failure Sending a Message, Closing Connection")
			return
		}
	}
}
