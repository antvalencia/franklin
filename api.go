package main

import (
	"franklin/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := models.SetupModels()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Run() // served on 0.0.0.0:8080
}
