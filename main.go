package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/eduncan911/podcast"
)

func main() {

	// feed, err := rss.Fetch("https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw")
	dat, err1 := ioutil.ReadFile("rss.yt.xml")

	feed, err2 := rss.Parse(dat)

	if err1 != nil || err2 != nil {
		fmt.Println(err1.Error(), err1.Error())
	}
	now := time.Now()

	p := podcast.New(feed.Title, feed.Link, feed.Description, &now, &feed.Refresh)
	p.Language = "ru-RU"

	// pItems := make([]*podcast.Item, 0, len(feed.Items))

	for _, ep := range feed.Items {
		item := new(podcast.Item)

		item.Title = ep.Title
		item.Link = ep.Link
		item.Description = ep.Desc
		item.PubDate = &ep.Date

		_, err := p.AddItem(*item)
		fmt.Println(err)
	}

	var buf bytes.Buffer
	err3 := p.Encode(&buf)
	res := buf.String() // s == "Size: 85 MB."
	fmt.Println(res, err3)

}
