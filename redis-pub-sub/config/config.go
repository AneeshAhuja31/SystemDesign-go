package config

import(
	"fmt"
	"context"
	"github.com/redis/go-redis/v9"
)

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error found: %s\n",err)
		panic(err)
	}
}

func ReturnContext() context.Context{
	return context.Background()
} 

func SetupRedis() *redis.Client{
	fmt.Printf("Setup redis service")
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
	})
}