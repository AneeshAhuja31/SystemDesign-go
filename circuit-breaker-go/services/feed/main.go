package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type post struct {
	ID int`json:"id"`
	Email string`json:"email"`
	Content string `json:"content"`
	Views int `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	godotenv.Load()
	router := gin.Default()

	router.GET("/feed", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "email not set in query param",
			})
			return
		}
		TRENDING_URL,ok := os.LookupEnv("TRENDING_URL")
		if !ok{
			TRENDING_URL= "http://localhost:7002"
		}
		RECOMMENDATIONS_URL,ok := os.LookupEnv("RECOMMENDATIONS_URL")
		if !ok {
			RECOMMENDATIONS_URL = "http://localhost:7003"
		}

		trendingResp, err := http.Get(TRENDING_URL + "/trending?limit=10")
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "failed to fetch trending posts: " + err.Error(),
			})
			return
		}
		defer trendingResp.Body.Close()

		var trendingPosts []post
		if err := json.NewDecoder(trendingResp.Body).Decode(&trendingPosts); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse trending posts: " + err.Error(),
			})
			return
		}

		recommendedResp, err := http.Get(RECOMMENDATIONS_URL + "/recommendations?email=" + email + "&limit=20")
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "failed to fetch recommendations: " + err.Error(),
			})
			return
		}
		defer recommendedResp.Body.Close()

		var recommendedPosts []post
		if err := json.NewDecoder(recommendedResp.Body).Decode(&recommendedPosts); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse recommendations: " + err.Error(),
			})
			return
		}

		seenIDs := make(map[int]bool)
		feedPosts := []post{}

		for _, p := range trendingPosts {
			if !seenIDs[p.ID] {
				seenIDs[p.ID] = true
				feedPosts = append(feedPosts, p)
			}
		}

		for _, p := range recommendedPosts {
			if !seenIDs[p.ID] {
				seenIDs[p.ID] = true
				feedPosts = append(feedPosts, p)
			}
		}

		ctx.JSON(http.StatusOK, feedPosts)
	})

	router.Run(":7004")
}
