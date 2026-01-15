package main

import (
	"time"
	"redis-pub-sub/config"
	"redis-pub-sub/publisher"
	"redis-pub-sub/subscriber"
)


func main(){
	rdb := config.SetupRedis()
	ctx := config.ReturnContext()
	go func(){
		subscriber.Subscriber(ctx,rdb,"orders")
	}()

	time.Sleep(1*time.Second)

	publisher.Publish(ctx,rdb,"orders","orders_created")
	publisher.Publish(ctx,rdb,"orders","order_shipped")
	time.Sleep(3 * time.Second)

}