import argparse
import json
import time
from datetime import datetime, timezone

import pika
import scrapy
from scrapy.crawler import CrawlerProcess


class CrawlWorkerSpider(scrapy.Spider):
    name = "crawl_worker"

    def __init__(self, tasks=None, rmq_channel=None, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.tasks = tasks or []
        self.rmq_channel = rmq_channel

    def start_requests(self):
        for task in self.tasks:
            yield scrapy.Request(
                url=task["url"],
                callback=self.parse_page,
                meta={"depth": task["depth"], "job_id": task["job_id"]},
                dont_filter=True,
            )

    def parse_page(self, response):
        links = [response.urljoin(href) for href in response.xpath("//a/@href").getall()]

        result = {
            "url": response.url,
            "status_code": response.status,
            "title": response.xpath("//title/text()").get(""),
            "body": " ".join(response.xpath("//body//text()").getall()).strip()[:50000],
            "links": links,
            "crawled_at": datetime.now(timezone.utc).isoformat(),
            "content_length": len(response.body),
            "depth": response.meta.get("depth", 0),
            "job_id": response.meta.get("job_id", ""),
        }

        self.rmq_channel.basic_publish(
            exchange="",
            routing_key="crawl_results",
            body=json.dumps(result),
            properties=pika.BasicProperties(content_type="application/json"),
        )
        self.logger.info(f"Published result for {response.url}")


def consume_batch(channel, queue_name, max_messages=50):
    tasks = []
    for _ in range(max_messages):
        method, _, body = channel.basic_get(queue=queue_name, auto_ack=True)
        if method is None:
            break
        task = json.loads(body)
        tasks.append(task)
    return tasks


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--worker-id", type=int, required=True, help="Worker ID (0-4)")
    parser.add_argument("--host", type=str, default="localhost", help="RabbitMQ host")
    args = parser.parse_args()

    queue_name = f"crawl_tasks_{args.worker_id}"
    print(f"[Worker {args.worker_id}] Consuming from {queue_name}")

    while True:
        connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=args.host)
        )
        channel = connection.channel()
        channel.queue_declare(queue=queue_name, durable=True)
        channel.queue_declare(queue="crawl_results", durable=True)

        tasks = consume_batch(channel, queue_name)

        if not tasks:
            connection.close()
            print(f"[Worker {args.worker_id}] No tasks, sleeping 5s...")
            time.sleep(5)
            continue

        print(f"[Worker {args.worker_id}] Got {len(tasks)} tasks, crawling...")

        process = CrawlerProcess(settings={
            "DUPEFILTER_CLASS": "scrapy.dupefilters.BaseDupeFilter",
            "ROBOTSTXT_OBEY": False,
            "DOWNLOAD_DELAY": 1,
            "USER_AGENT": "WebCrawler/1.0",
            "LOG_LEVEL": "INFO",
            "CONCURRENT_REQUESTS": 8,
        })
        process.crawl(CrawlWorkerSpider, tasks=tasks, rmq_channel=channel)
        process.start()

        connection.close()
        print(f"[Worker {args.worker_id}] Batch done, polling again...")


if __name__ == "__main__":
    main()
