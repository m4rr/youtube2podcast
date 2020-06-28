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

var episodes []itcEpisode

func main() {

	db, dbErr := gorm.Open("sqlite3", "./sqlite3.db")
	defer db.Close()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	dat, readErr := ioutil.ReadFile("rss.yt.xml")
	feed, parsErr := rss.Parse(dat)
	if readErr != nil || parsErr != nil {
		fmt.Println(readErr.Error(), parsErr.Error())
	}

	// feed, err2 := rss.Fetch("https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw")
	// if err2 != nil {
	// 	fmt.Println("err 2", err2.Error())
	// }

	now := time.Now()

	p := podcast.New(feed.Title, feed.Link, feed.Description, &now, &feed.Refresh)
	p.Language = "ru-RU"

	// pItems := make([]*podcast.Item, 0, len(feed.Items))

	for _, ytEpisode := range feed.Items {
		itcItem := new(podcast.Item)

		itcItem.Title = ytEpisode.Title
		itcItem.PubDate = &ytEpisode.Date

		itcItem.Link = ytEpisode.Link
		itcItem.Description = ytEpisode.Desc

		if len(ytEpisode.Desc) == 0 {
			itcItem.Description = "<No Shownotes>"
		}

		author := podcast.Author{}
		author.Name = feed.Title
		itcItem.Author = &author

		itcItem.Comments = strconv.Itoa(ytEpisode.Views) + " Views"

		_, addErr := p.AddItem(*itcItem)
		if addErr != nil {
			fmt.Println(addErr)
		}

		episodes = append(episodes, itcEpisode{gorm.Model{}, *itcItem})
	}

	writErr := ioutil.WriteFile("rss.ktotoneprav.xml", p.Bytes(), 0644)
	if writErr != nil {
		log.Fatal(writErr)
	}

	itcPodc := itcPodcast{gorm.Model{}, p, episodes}

	db.AutoMigrate(&itcEpisode{}, &itcPodcast{})
	db.Create(&itcPodc)

}

type itcPodcast struct {
	gorm.Model
	podcast.Podcast
	episodes []itcEpisode
}

type itcEpisode struct {
	gorm.Model
	podcast.Item
}
