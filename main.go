package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	db, dbErr := gorm.Open("sqlite3", "./sqlite3.db")
	defer db.Close()
	// db = db.Debug()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	ktotonepravURL := "https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw"
	feed, parsErr := readRSS(&ktotonepravURL)
	if parsErr != nil {
		log.Fatal(parsErr.Error())
	}

	db.AutoMigrate(&Author{}, &Episode{}, &Category{}, &Podcast{})

	for _, cat := range fillCategories() {
		if db.NewRecord(cat) {
			db.Create(&cat)
		}
	}

	parsedPod := parseYtRss(feed)
	db.Create(&parsedPod)

	var selectedPod Podcast
	db.Preload("Episodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("Categories").Last(&selectedPod)

	itcPod := itcPodcastFrom(&selectedPod)
	writErr := writeItunesPodcastRssXML(itcPod)
	if writErr != nil {
		log.Fatal(writErr)
	}

	runWebServer(selectedPod)
}
