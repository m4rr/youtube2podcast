package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rss "github.com/m4rr/yt-rss"
)

type TheAuthor struct {
	gorm.Model

	Name string
	Link string
}

type TheCategory struct {
	gorm.Model

	Name string
}

type ThePodcast struct {
	gorm.Model

	// YtID           string `xml:"id"`
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	AuthorName     string    `xml:"author.name"`
	FirstPublished time.Time `xml:"published"`
	Description    string

	TheEpisodes []TheEpisode

	Lang        string
	TheCategory []TheCategory
	Cached      time.Time
}

type TheEpisode struct {
	gorm.Model

	ThePodcast   ThePodcast
	ThePodcastID uint

	YtID           string    `xml:"id"`
	VideoID        string    `xml:"videoId"`
	ChannelID      string    `xml:"channelId"`
	Title          string    `xml:"title"`
	YtLink         string    `xml:"link,href,attr"`
	Author         TheAuthor `xml:"author"`
	Published      time.Time `xml:"published"`
	Updated        string    `xml:"updated"`
	CoverImageLink string    `xml:"group.thumbnail.url"`
	Description    string    `xml:"group.description"`
	StarRating     float64   `xml:"community.starRating"`
	Views          int       `xml:"community.statistics"`
}

func parseYtRss(feed *rss.Feed) (thePod ThePodcast) {

	thePod.Lang = "ru-RU"
	thePod.Title = feed.Title
	thePod.Link = feed.Link
	thePod.AuthorName = feed.Nickname

	theEps := []TheEpisode{}

	for _, ytEpisode := range feed.Items {
		theEp := TheEpisode{}
		theEp.Title = ytEpisode.Title
		theEp.Published = ytEpisode.Date
		theEp.YtLink = ytEpisode.Link

		theEp.Description = ytEpisode.Desc
		if len(ytEpisode.Desc) == 0 {
			theEp.Description = "<No Shownotes>"
		}

		author := TheAuthor{}
		author.Name = ytEpisode.Author
		theEp.Author = author

		theEp.Views = ytEpisode.Views

		theEps = append(theEps, theEp)
	}

	thePod.TheEpisodes = theEps
	thePod.Cached = time.Now()

	return
}

func itcPodcastFrom(thePod *ThePodcast) (p podcast.Podcast) {

	p = podcast.New(thePod.Title, thePod.Link, thePod.Description, &thePod.FirstPublished, &thePod.UpdatedAt)
	p.Language = "ru-RU"

	// pItems := make([]*podcast.Item, 0, len(feed.Items))

	for _, ytEpisode := range thePod.TheEpisodes {
		itcItem := new(podcast.Item)

		itcItem.Title = ytEpisode.Title
		itcItem.PubDate = &ytEpisode.ThePodcast.Cached

		itcItem.Link = ytEpisode.YtLink
		itcItem.Description = ytEpisode.Description

		author := podcast.Author{}
		author.Name = thePod.AuthorName
		itcItem.Author = &author

		itcItem.Comments = strconv.Itoa(ytEpisode.Views) + " Views"

		_, addErr := p.AddItem(*itcItem)
		if addErr != nil {
			fmt.Println(addErr)
		}
	}

	return
}

func main() {

	db, dbErr := gorm.Open("sqlite3", "./sqlite3.db")
	db = db.Debug()
	defer db.Close()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	db.AutoMigrate(&TheAuthor{}, &TheEpisode{}, &ThePodcast{})

	// p := podcast.Item{}

	dat, readErr := ioutil.ReadFile("rss.yt.xml")
	feed, parsErr := rss.Parse(dat)
	if readErr != nil || parsErr != nil {
		fmt.Println(readErr.Error(), parsErr.Error())
	}

	feed, err2 := rss.Fetch("https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw")
	if err2 != nil {
		fmt.Println("err 2", err2.Error())
	}

	thePod := parseYtRss(feed)

	db.Create(&thePod)

	var thePod2 ThePodcast

	db.Preload("TheEpisodes").Last(&thePod2)

	itcPodcast := itcPodcastFrom(&thePod2)

	writErr := ioutil.WriteFile("rss.itc.generated.xml", itcPodcast.Bytes(), 0644)
	if writErr != nil {
		log.Fatal(writErr)
	}
}
