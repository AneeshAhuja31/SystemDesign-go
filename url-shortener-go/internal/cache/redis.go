package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)


func InitRedis(host string, port int)(*redis.Client,context.Context){
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + fmt.Sprint(port),
		Password: "",
		DB: 0,
	})
	return rdb,ctx
}

