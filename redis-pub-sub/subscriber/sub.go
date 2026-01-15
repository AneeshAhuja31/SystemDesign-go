package subscriber

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"redis-pub-sub/config"
)

func Subscriber(ctx context.Context, rdb *redis.Client, channel string) error{
	sub := rdb.Subscribe(ctx, channel)
	defer func (){
		closeErr := sub.Close()
		config.HandleError(closeErr)
	}()
	ch := sub.Channel()
	fmt.Println("Subscribed to:", channel)
	for msg := range ch{
		fmt.Println("Received message: ", msg.Payload)
	}
	return nil
}