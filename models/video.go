package models

import "time"

type Video struct {
	ID           string    `json:"id" bson:"_id"`
	Title        string    `json:"title" bson:"title"`
	URL          string    `json:"url" bson:"url"`
	ThumbnailURL string    `json:"thumbnail_url" bson:"thumbnail_url"`
	Description  string    `json:"description" bson:"description"`
	Category     string    `json:"category" bson:"category"`
	Duration     int       `json:"duration" bson:"duration"`
	Views        int       `json:"views" bson:"views"`
	UploaderID   string    `json:"uploader_id" bson:"uploader_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func NewVideo(title string, url string, thumbnailURL string, description string, category string, duration int, uploaderID string) *Video {
	return &Video{
		Title:        title,
		URL:          url,
		ThumbnailURL: thumbnailURL,
		Description:  description,
		Category:     category,
		Duration:     duration,
		Views:        0,
		UploaderID:   uploaderID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (v *Video) SetID(id string) {
	v.ID = id
}
