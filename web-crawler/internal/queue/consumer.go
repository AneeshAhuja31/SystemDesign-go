package queue

import (
	"encoding/json"
	"web-crawler/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)
func ConsumeCrawlResults(ch *amqp.Channel, handler func(models.CrawlResult)){
	msgs,_ := ch.Consume("crawl_results","",true,false,false,false,nil)
	for msg := range msgs{
		var results models.CrawlResult
		json.Unmarshal(msg.Body,&results)
		handler(results)
	}
}