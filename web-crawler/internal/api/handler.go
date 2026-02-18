package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"web-crawler/internal/bloom"
	"web-crawler/internal/db"
	"web-crawler/internal/models"
)

type API struct {
	DB    *sql.DB
	Bloom *bloom.BloomFilter
}

func NewAPI(pg *sql.DB, bf *bloom.BloomFilter) *API {
	return &API{DB: pg, Bloom: bf}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/seed", a.AddSeedURLs)
	mux.HandleFunc("POST /api/reindex", a.TriggerReindex)
	mux.HandleFunc("GET /api/stats", a.GetCrawlStats)
	mux.HandleFunc("GET /api/frontier", a.GetFrontierStatus)
}

func (a *API) AddSeedURLs(w http.ResponseWriter, r *http.Request) {
	var req models.SeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	added := 0
	skipped := 0

	for _, rawURL := range req.URLs {
		parsed, err := url.Parse(rawURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			skipped++
			continue
		}

		if a.Bloom.MightContain(rawURL) {
			skipped++
			continue
		}

		domain := parsed.Hostname()
		err = db.EnqueueURL(a.DB, rawURL, domain, 0, 1)
		if err != nil {
			log.Println("Error enqueuing seed URL: ", err)
			skipped++
			continue
		}

		a.Bloom.Add(rawURL)
		added++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"added":              added,
		"skipped_duplicates": skipped,
	})
}

func (a *API) TriggerReindex(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(req.URL)
	if err != nil || parsed.Scheme == "" {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	domain := parsed.Hostname()
	err = db.EnqueueURL(a.DB, req.URL, domain, 0, 1)
	if err != nil {
		log.Println("Error enqueuing reindex URL: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "reindex_requested",
	})
}

func (a *API) GetCrawlStats(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetCrawlStats(a.DB)
	if err != nil {
		http.Error(w, "error fetching stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (a *API) GetFrontierStatus(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetCrawlStats(a.DB)
	if err != nil {
		http.Error(w, "error fetching frontier status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"frontier": stats,
	})
}
