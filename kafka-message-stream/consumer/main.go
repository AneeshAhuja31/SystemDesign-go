package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func handleErr(err error){
	if err != nil{
		fmt.Println(err)
		panic(err)
	}
}

func main(){
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":"localhost:9092",
		"group.id": "order-group",
		"auto.offset.reset":"earliest",
	})
	handleErr(err)
	defer c.Close()

	c.SubscribeTopics([]string{"orders"}, nil)

	sigs := make(chan os.Signal,1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case <- sigs:
			run = false
		default: 
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch msg := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Consumed: %s at %v\n", string(msg.Value), msg.TopicPartition)
			case kafka.Error:
				fmt.Printf("Error: %v\n", msg)
			}
		}
	}
	fmt.Println("Consumer closing")
}