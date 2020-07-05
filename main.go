package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/eduncan911/podcast"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rss "github.com/m4rr/yt-rss"
)

func writeItunesPodcastRssXML(itcPodcast podcast.Podcast) {

	writErr := ioutil.WriteFile("rsss/rss.itc.gen.xml", itcPodcast.Bytes(), 0644)
	if writErr != nil {
		log.Fatal(writErr)
	}

}

func runWebServer(tehPod Podcast) {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", tehPod)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

var db *gorm.DB

func main() {

	db, _ = gorm.Open("sqlite3", "./sqlite3.db")
	defer db.Close()
	// db = db.Debug()
	// if dbErr != nil {
	// 	log.Fatal(dbErr)
	// }

	dat, readErr := ioutil.ReadFile("rsss/rss.yt.orig.xml")
	feed, parsErr := rss.Parse(dat)
	if readErr != nil || parsErr != nil {
		log.Fatal(readErr.Error(), parsErr.Error())
	}

	// feed, err2 := rss.Fetch("https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw")
	// if err2 != nil {
	// 	log.Fatal("err 2", err2.Error())
	// }

	db.AutoMigrate(&Author{}, &Episode{}, &Category{}, &Podcast{})

	thePod := parseYtRss(feed)
	db.Create(&thePod)

	var thePod2 Podcast
	db.Preload("Episodes").Preload("Categories").Last(&thePod2)
	sort.Sort(ByID(thePod2.Episodes))

	itcPodcast := itcPodcastFrom(&thePod2)
	writeItunesPodcastRssXML(itcPodcast)

	// runWebServer(thePod2)
	writeItunesPodcastRssXML(itcPodcast)

}
