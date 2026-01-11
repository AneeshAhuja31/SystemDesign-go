package main

import (
	"fmt"


	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func handleErr(err error){
	if err != nil{
		fmt.Println(err)
		panic(err)
	}
}

func main(){
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers":"localhost:9092"})
	handleErr(err)

	defer p.Close()

	go func(){
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "orders"
	values := []string{"order-1","order-2","order-3"}
	for _,v := range values{
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value: []byte(v),
		}, nil)
	}  
	p.Flush(15000)
	fmt.Println("All messages sent!")
}