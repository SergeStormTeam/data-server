package redis

import (
	"context"
)

func PublishToDashboard(payload any) error {
	ctx := context.Background()

	err := chanClient.Publish(ctx, "zephyr-update", payload).Err()
	if err != nil {
		return err
	}

	return nil
}
