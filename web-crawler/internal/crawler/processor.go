package crawler

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/elastic/go-elasticsearch/v9"
	amqp "github.com/rabbitmq/amqp091-go"
	"web-crawler/internal/bloom"
	"web-crawler/internal/db"
	"web-crawler/internal/elastic"
	"web-crawler/internal/models"
	"web-crawler/internal/queue"
)

type Processor struct {
	DB *sql.DB
	Ch *amqp.Channel
	ES *elasticsearch.Client
	Bloom *bloom.BloomFilter
}

func NewProcessor(pg *sql.DB, ch *amqp.Channel, es *elasticsearch.Client, bf *bloom.BloomFilter) *Processor {
	return &Processor{
		DB: pg,
		Ch: ch,
		ES: es,
		Bloom: bf,
	}
}

func (p *Processor) Start() {
	log.Println("Processor started, consuming crawl_results...")

	queue.ConsumeCrawlResults(p.Ch, func(result models.CrawlResult) {
		doc := models.PageDocument{
			URL: result.URL,
			StatusCode:    result.StatusCode,
			Title: result.Title,
			Body: result.Body,
			Links: result.Links,
			CrawledAt: result.CrawledAt,
			ContentLength: result.ContentLength,
			Depth: result.Depth,
			JobID: result.JobID,
		}

		//extract domain from url
		parsed, err := url.Parse(result.URL)
		if err == nil {
			doc.Domain = parsed.Hostname()
		}

		err = elastic.IndexPage(p.ES, doc)
		if err != nil {
			log.Println("Error indexing page: ", err)
		} else {
			log.Printf("Indexed %s in Elasticsearch", result.URL)
		}

		//discover new urls from links
		for _, link := range result.Links {
			parsedLink, err := url.Parse(link)
			if err != nil || parsedLink.Scheme == "" || parsedLink.Host == "" {
				continue
			}
			normalized := parsedLink.Scheme + "://" + parsedLink.Host + parsedLink.Path

			if p.Bloom.MightContain(normalized) {
				continue
			}

			p.Bloom.Add(normalized)
			domain := parsedLink.Hostname()
			err = db.EnqueueURL(p.DB, normalized, domain, result.Depth+1, 5)
			if err != nil {
				log.Println("Error enqueuing discovered URL: ", err)
			}
		}

		//mark original url as crawled in frontier
		//extract ID from job_id (format: "job-{id}")
		var urlID int
		fmt.Sscanf(result.JobID, "job-%d", &urlID)
		if urlID > 0 {
			db.MarkCrawled(p.DB, urlID, result.StatusCode)
		}
	})
}
