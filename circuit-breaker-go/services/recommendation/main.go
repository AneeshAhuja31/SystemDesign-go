package main

import (
	"os"
	"time"
	"fmt"
	"net/http"
	"io"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	
)

func main(){
	router := gin.Default()
	router.GET("/recommendations",func(ctx *gin.Context) {
		email := ctx.Query("email")
		
	})
}