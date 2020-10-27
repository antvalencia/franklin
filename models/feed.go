package models

import "time"

type Feed struct {
	ID                   uint      `json:"id" gorm:"primary_key"`
	URL                  string    `json:"url"`
	Category             string    `json:"category"`
	Source               string    `json:"source"`
	Title                string    `json:"title"`
	DownloadFrequency    int       `json:"download_frequency" gorm:"default:1"` // download every 1 day
	PreviousDownloadTime time.Time `json:"previous_download_time"`
}

type UpdateFeedInput struct {
	ID                   uint      `json:"id" gorm:"primary_key"`
	PreviousDownloadTime time.Time `json:"previous_download_time"`
}

type FeedsGetResponse struct {
	Data []Feed `json:"data"`
}
