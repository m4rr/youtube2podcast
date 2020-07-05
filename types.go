package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Author struct {
	gorm.Model

	Name  string
	Link  string
	Email string
}

type Category struct {
	// gorm.Model

	ID          uint       `gorm:"primary_key,auto_increment"`
	Name        string     `gorm:"unique;not_null"`
	Podcasts    []*Podcast `gorm:"many2many:podcast_categories;"`
	PodcastYtID string
}

type Episode struct {
	// gorm.Model

	ID          uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	YtID        string `gorm:"UNIQUE_INDEX;NOT_NULL"`
	PodcastYtID string
	Podcast     Podcast

	VideoID        string    `xml:"videoId"`
	ChannelID      string    `xml:"channelId"`
	Title          string    `xml:"title"`
	YtLink         string    `xml:"link,href,attr"`
	Author         Author    `xml:"author"`
	Published      time.Time `xml:"published"`
	Updated        string    `xml:"updated"`
	CoverImageLink string    `xml:"group.thumbnail.url"`
	Description    string    `xml:"group.description"`
	StarRating     float64   `xml:"community.starRating"`
	Views          int       `xml:"community.statistics"`
}

type Podcast struct {
	// gorm.Model

	ID         uint   `gorm:"AUTO_INCREMENT"`
	YtID       string `gorm:"PRIMARY_KEY;UNIQUE_INDEX;NOT_NULL"`
	Episodes   []Episode
	Categories []*Category `gorm:"many2many:podcast_categories;"`

	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	AuthorName     string    `xml:"author.name"`
	FirstPublished time.Time `xml:"published"`
	Description    string

	Lang   string
	Cached time.Time
}

type ByID []Episode

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
