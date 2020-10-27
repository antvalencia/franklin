package resources

import (
	"net/http"

	models "franklin/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetFeeds(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var feeds []models.Feed
	db.Find(&feeds)

	c.JSON(http.StatusOK, gin.H{"data": feeds})
}

func PatchFeed(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var feed models.Feed
	if err := db.Where("id = ?", c.Param("id")).First(&feed).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input models.UpdateFeedInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&feed).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": feed})
}
