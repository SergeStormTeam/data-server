package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SergeStormTeam/dashboard-handler/dashboard"
	"github.com/SergeStormTeam/dashboard-handler/logging"
	"github.com/SergeStormTeam/dashboard-handler/types"
	"github.com/sirupsen/logrus"
)

func SubscribeToZephyrApi() error {
	ctx := context.Background()

	sub := chanClient.Subscribe(ctx, "zephyr-update")
	defer sub.Close()

	subChannel := sub.Channel()

	for {
		select {
		case payload, ok := <-subChannel:
			if !ok {
				return fmt.Errorf("redis subscription closed")
			}

			var update types.ZephyrUpdate
			err := json.Unmarshal([]byte(payload.Payload), &update)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "redis", "method": "InitRedis"}).Warn("Error unmarshaling message!")
				continue
			}

			dashboard.MessageAllWebsockets(update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
