package queue

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal("Error initializing RabbitMQ: ", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error opening RabbitMQ channel: ", err)
	}
	return conn, ch
}

func DeclareQueues(ch *amqp.Channel) {
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("crawl_tasks_%d", i)
		_, err := ch.QueueDeclare(name, true, false, false, false, nil)
		if err != nil {
			log.Fatal("Error declaring queue ", name, ": ", err)
		}
	}
	_, err := ch.QueueDeclare("crawl_results", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error declaring crawl_results queue: ", err)
	}
	log.Println("All RabbitMQ queues declared")
}