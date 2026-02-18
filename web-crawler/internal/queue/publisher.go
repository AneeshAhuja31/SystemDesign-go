package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"web-crawler/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)
func PublishCrawlTask(ch *amqp.Channel, workerID int, task models.CrawlTask) error {
	data,err:= json.Marshal(task)
	if err != nil{
		log.Println("Error converting to json format: ",err)
		return err
	}
	ctx := context.Background()
	err = ch.PublishWithContext(ctx,"","crawl_tasks_"+fmt.Sprint(workerID),false,false,amqp.Publishing{
		ContentType: "application/json",
		Body: data,
	})
	return err
}