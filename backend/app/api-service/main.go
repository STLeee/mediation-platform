package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
)

func init() {
	// Load config
	_, err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config.GetConfig()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
