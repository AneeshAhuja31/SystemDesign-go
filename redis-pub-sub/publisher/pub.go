package publisher

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"redis-pub-sub/config"
)

func Publish(ctx context.Context,rdb *redis.Client,channel string, message string) error{
	err := rdb.Publish(ctx,channel,message).Err()
	config.HandleError(err)
	fmt.Println("Published: ",message)
	return nil
}