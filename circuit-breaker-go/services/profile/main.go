package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main(){
	db := initPQ("localhost",5432)
	router := gin.Default()
	router.GET("/profile", func(c*gin.Context){
		email := c.Query("email")
		if email == "" {
			c.JSON(400, gin.H{
				"error":"Empty email",
			})
			return
		}
		userProfile,err := fetchProfileData(db,email)
		if err != nil {
			fmt.Println("Error: ",err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,userProfile)
	})
	err := router.Run(":7000")
	handleError(err)
}
