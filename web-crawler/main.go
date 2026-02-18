package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"web-crawler/internal/api"
	"web-crawler/internal/bloom"
	"web-crawler/internal/crawler"
	"web-crawler/internal/db"
	"web-crawler/internal/elastic"
	"web-crawler/internal/hashring"
	"web-crawler/internal/queue"
	"web-crawler/internal/scheduler"
)

func main() {
	godotenv.Load()

	// Initialize PostgreSQL
	pg := db.InitDB()
	log.Println("PostgreSQL connected")

	// Initialize RabbitMQ
	conn, ch := queue.InitRabbitMQ()
	defer conn.Close()
	defer ch.Close()
	queue.DeclareQueues(ch)
	log.Println("RabbitMQ connected")

	// Initialize Elasticsearch
	es := elastic.InitClient()
	log.Println("Elasticsearch connected")

	// Initialize Redis + Hash Ring (5 workers, 10 vnodes each)
	rdb, ctx := hashring.InitRedis()
	ring := hashring.NewHashRing(rdb, ctx, 5, 10)
	log.Println("Redis connected, hash ring initialized")

	// Initialize Bloom Filter (1M bits, 5 hash functions)
	bf := bloom.NewBloomFilter(1_000_000, 5)
	log.Println("Bloom filter initialized")

	// Initialize Politeness Enforcer (1s min delay per domain)
	politeness := scheduler.NewPolitenessEnforcer(1 * time.Second)

	// Start Scheduler goroutine
	sched := scheduler.NewScheduler(pg, ch, bf, ring, politeness)
	go sched.Start()

	// Start Processor goroutine
	proc := crawler.NewProcessor(pg, ch, es, bf)
	go proc.Start()

	// Start HTTP API
	apiServer := api.NewAPI(pg, bf)
	mux := http.NewServeMux()
	apiServer.RegisterRoutes(mux)

	log.Println("API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
