package main

import (
	"franklin/models"

	resources "franklin/resources"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := models.SetupModels()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/feeds", resources.GetFeeds)
	r.PATCH("/feeds/:id", resources.PatchFeed)
	r.GET("/news_items", resources.GetNewsItems)
	r.GET("/news_items/:guid_md5", resources.GetNewsItem)
	r.POST("/news_items", resources.PostNewsItem)

	r.Run()
}
