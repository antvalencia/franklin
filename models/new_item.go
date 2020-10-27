package models

import "time"

type NewsItem struct {
	ID      uint `gorm:"primary_key"`
	Title   string
	Link    string
	Desc    string
	Guid    string
	GuidMD5 string
	PubDate time.Time
	Feed    Feed
	FeedID  uint
}

type CreateNewsItemInput struct {
	Title   string    `json:"title" binding:"required"`
	Link    string    `json:"link" binding:"required"`
	Desc    string    `json:"description"`
	Guid    string    `json:"guid" binding:"required"`
	GuidMD5 string    `json:"guid_md5" binding:"required"`
	PubDate time.Time `json:"pub_date"`
	FeedID  uint      `json:"feed_id,string" binding:"required"`
}

type NewsItemGetResponse struct {
	Data NewsItem `json:"data"`
}
