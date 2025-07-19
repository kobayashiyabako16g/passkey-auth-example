package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/handler"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/server"
)

func enableCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(enableCORSMiddleware())

	api := router.Group("/api")
	{
		api.GET("/register/begin", handler.BeginRegistration(srv))
		// 同様に他のルートも登録
	}

	fmt.Println("Gin API server starting on :8080")
	router.Run(":8080")
}
