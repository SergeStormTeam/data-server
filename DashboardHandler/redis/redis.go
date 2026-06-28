package redis

import (
	"context"
	"os"

	"github.com/SergeStormTeam/dashboard-handler/logging"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

var client *redis.Client
var chanClient *redis.Client

func InitRedis() error {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	newClient := redis.NewClient(opt)
	_, err = newClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	newChanClient := redis.NewClient(opt)
	_, err = newChanClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	client = newClient
	chanClient = newChanClient
	logging.Logger.WithFields(logrus.Fields{"module": "redis", "method": "InitRedis"}).Info("Successfully Initialized Redis!")
	return nil
}
