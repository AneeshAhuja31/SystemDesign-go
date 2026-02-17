package db

import (
	"database/sql"
	"log"
	"web-crawler/internal/models"
)
func EnqueueURL(db *sql.DB, url string, domain string, depth int, priority int)error{
	enqueueString := "INSERT INTO frontier (url, domain, depth, priority) VALUES ($1, $2, $3, $4) ON CONFLICT (url) DO NOTHING"
	_,err := db.Exec(enqueueString,url,domain,depth,priority)
	if err != nil {
		log.Println("Error enqueuing url, ",err)
		return err
	}
	return nil
}

func DequeueURLs(db *sql.DB, limit int)([]models.FrontierURL,error){
	queryString := `
		SELECT id, url, domain, depth 
		FROM frontier 
		WHERE status='pending' ORDER BY priority, created_at LIMIT $1	
	`
	urlsToDequeue,err := db.Query(queryString,limit)
	if err != nil {
		log.Println("Error in postgres dequeue select query: ",err)
		return []models.FrontierURL{},err
	}
	defer urlsToDequeue.Close()
	var frontierUrls []models.FrontierURL
	for urlsToDequeue.Next(){
		var frontierUrl models.FrontierURL
		urlsToDequeue.Scan(&frontierUrl.ID,&frontierUrl.URL,&frontierUrl.Domain,&frontierUrl.Depth)
		frontierUrls = append(frontierUrls, frontierUrl)
	}
	for _,frontierUrl := range(frontierUrls){
		updateExecString := "UPDATE frontier SET status = 'in_progress' WHERE id = $1"
		_,err := db.Exec(updateExecString,frontierUrl.ID)
		if err != nil {
			log.Println("Error updating frontieer urls to in_progress: ",err)
			return []models.FrontierURL{},err
		}
	}
	return frontierUrls,nil
}

func MarkCrawled(db *sql.DB, id int, statusCode int) error {
	_,err := db.Exec("UPDATE frontier SET status='crawled', status_code=$1, updated_at=NOW() WHERE id=$2",statusCode,id)
	return err
}

func MarkFailed(db *sql.DB, id int, reason string) error {
	_,err := db.Exec("UPDATE frontier SET status='failed', error_msg=$1, updated_at=NOW() WHERE id=$2",reason,id)
	return err
}

func GetCrawlStats(db *sql.DB) (*models.CrawlStats,error){
	crawlstatrows := db.QueryRow(`SELECT 
				COUNT(*) as total, 
				COUNT(*) FILTER (WHERE status='pending'), 
				COUNT(*) FILTER (WHERE status='in_progress'),
				COUNT(*) FILTER (WHERE status='crawled'),
				COUNT(*) FILTER (WHERE status='failed')
			`)
	
	var crawlstats models.CrawlStats
	crawlstatrows.Scan(&crawlstats.TotalURLs,&crawlstats.Pending,&crawlstats.InProgress,&crawlstats.Crawled,&crawlstats.Failed)
	return &crawlstats,nil
}
