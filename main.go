package main

import (
	"log"
	"sort"

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

	feed, parsErr := readRSS(nil)
	if parsErr != nil {
		log.Fatal(parsErr.Error())
	}

	db.AutoMigrate(&Author{}, &Episode{}, &Category{}, &Podcast{})

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
