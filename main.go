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

	thePod := parseYtRss(feed)
	db.Create(&thePod)

	var thePod2 Podcast
	db.Preload("Episodes").Preload("Categories").Last(&thePod2)
	sort.Sort(ByID(thePod2.Episodes))

	itcPodcast := itcPodcastFrom(&thePod2)
	writErr := writeItunesPodcastRssXML(itcPodcast)
	if writErr != nil {
		log.Fatal(writErr)
	}

	runWebServer(thePod2)
}
