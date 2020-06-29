package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gin-gonic/gin"
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
	// gorm.Model

	ID          uint   `gorm:"AUTO_INCREMENT"`
	YtID        string `gorm:"PRIMARY_KEY;UNIQUE_INDEX;NOT_NULL"`
	TheEpisodes []TheEpisode

	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	AuthorName     string    `xml:"author.name"`
	FirstPublished time.Time `xml:"published"`
	Description    string

	Lang        string
	TheCategory []TheCategory
	Cached      time.Time
}

type TheEpisode struct {
	// gorm.Model

	ID             uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	YtID           string `gorm:"UNIQUE_INDEX;NOT_NULL"`
	ThePodcast     ThePodcast
	ThePodcastYtID string

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

func parseYtRss(feed *rss.Feed) ThePodcast {

	thePod := ThePodcast{}

	thePod.YtID = feed.ID
	thePod.Lang = "ru-RU"
	thePod.Title = feed.Title
	thePod.Link = feed.Link
	thePod.AuthorName = feed.Nickname

	theEps := []TheEpisode{}

	for _, ytEpisode := range feed.Items {
		theEp := TheEpisode{}

		theEp.YtID = ytEpisode.ID
		theEp.Title = ytEpisode.Title
		theEp.Published = ytEpisode.Date
		theEp.YtLink = ytEpisode.Link

		if len(ytEpisode.Desc) == 0 {
			theEp.Description = "<No Shownotes>"
		} else {
			theEp.Description = ytEpisode.Desc
		}

		author := TheAuthor{}
		author.Name = ytEpisode.Author
		theEp.Author = author

		theEp.Views = ytEpisode.Views

		theEp.ThePodcastYtID = thePod.YtID
		theEps = append(theEps, theEp)
	}

	thePod.TheEpisodes = theEps
	thePod.Cached = time.Now()

	return thePod
}

func itcPodcastFrom(thePod *ThePodcast) podcast.Podcast {

	p := podcast.New(thePod.Title, thePod.Link, thePod.Description, &thePod.FirstPublished, &thePod.Cached)

	p.IAuthor = thePod.Title //AuthorName
	p.Language = "ru-RU"
	p.IExplicit = "true"

	for _, ytEpisode := range thePod.TheEpisodes {
		itcItem := new(podcast.Item)

		itcItem.Title = ytEpisode.Title
		itcItem.PubDate = &ytEpisode.Published

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

	return p
}

func writeItunesPodcastRssXML(itcPodcast podcast.Podcast) {

	writErr := ioutil.WriteFile("rss.itc.generated.xml", itcPodcast.Bytes(), 0644)
	if writErr != nil {
		log.Fatal(writErr)
	}
}

func runWebServer(tehPod ThePodcast) {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", tehPod)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

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
	sort.Sort(ByID(thePod2.TheEpisodes))

	itcPodcast := itcPodcastFrom(&thePod2)
	writeItunesPodcastRssXML(itcPodcast)

	runWebServer(thePod2)

}

type ByID []TheEpisode

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
