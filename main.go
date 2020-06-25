package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	rss "github.com/m4rr/yt-rss"
)

func main() {

	dat, err1 := ioutil.ReadFile("rss.yt.xml")
	feed, err2 := rss.Parse(dat)
	if err1 != nil || err2 != nil {
		// fmt.Println(err1.Error(), err2.Error())
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
		itcItem.Link = ytEpisode.Link
		itcItem.Description = ytEpisode.Desc
		itcItem.PubDate = &ytEpisode.Date

		author := podcast.Author{}
		author.Name = feed.Title
		itcItem.Author = &author

		itcItem.Comments = strconv.Itoa(ytEpisode.Views) + " Views"

		_, err := p.AddItem(*itcItem)
		fmt.Println(err)
	}

	var buf bytes.Buffer
	err3 := p.Encode(&buf)
	res := buf.String() // s == "Size: 85 MB."
	fmt.Println(res, err3)

}
