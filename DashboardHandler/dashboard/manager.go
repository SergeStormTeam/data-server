package dashboard

import (
	"sync"

	"github.com/SergeStormTeam/dashboard-handler/logging"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const CHANNEL_SIZE int = 64

type ServerHub struct {
	Clients map[*websocket.Conn]*Client
	Mu      sync.RWMutex
}

var Hub ServerHub = ServerHub{
	Clients: make(map[*websocket.Conn]*Client),
	Mu:      sync.RWMutex{},
}

func AddWebSocketConnection(conn *websocket.Conn) {
	newClient := &Client{
		conn: conn,
		send: make(chan any, CHANNEL_SIZE),
	}

	Hub.Mu.Lock()
	logging.Logger.WithFields(logrus.Fields{"module": "dashboard", "method": "AddWebSocketConnection"}).Info("NEW CONNECTION!")
	Hub.Clients[conn] = newClient
	Hub.Mu.Unlock()

	go newClient.SendLoop()
}

func RemoveWebSocketConnection(conn *websocket.Conn) {
	Hub.Mu.Lock()
	c, ok := Hub.Clients[conn]
	if ok {
		delete(Hub.Clients, conn)
	}
	Hub.Mu.Unlock()

	if ok {
		c.Close()
	}
}

func MessageAllWebsockets(payload any) {
	Hub.Mu.RLock()
	defer Hub.Mu.RUnlock()

	logging.Logger.WithFields(logrus.Fields{"payload": payload, "module": "dashboard", "method": "AddWebSocketConnection"}).Info("New Payload")

	for conn, client := range Hub.Clients {
		select {
		case client.send <- payload:
		default:
			logging.Logger.WithFields(logrus.Fields{"module": "dashboard", "method": "MessageAllWebsockets"}).Warn("Websocket Connection Full, Closing")
			go RemoveWebSocketConnection(conn)
		}
	}
}
