package resources

import (
	"net/http"

	models "franklin/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetNewsItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var news_items []models.NewsItem
	if err := db.Find(&news_items).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
	}

	for i, _ := range news_items {
		db.Model(news_items[i]).Related(&news_items[i].Feed)
	}

	c.JSON(http.StatusOK, gin.H{"data": news_items})
}

func GetNewsItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var news_item models.NewsItem
	if err := db.Where("guid_md5 = ?", c.Param("guid_md5")).First(&news_item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": news_item})
	}
}

func PostNewsItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.CreateNewsItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news_item := models.NewsItem{
		Title:   input.Title,
		Link:    input.Link,
		Desc:    input.Desc,
		Guid:    input.Guid,
		GuidMD5: input.GuidMD5,
		PubDate: input.PubDate,
		FeedID:  input.FeedID,
	}
	db.Save(&news_item)

	c.JSON(http.StatusCreated, gin.H{"data": news_item})
}
