package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	rss "github.com/m4rr/yt-rss"
)

func readRSS(url *string) (feed *rss.Feed, parsErr error) {

	if url == nil {
		data, readErr := ioutil.ReadFile("rsss/rss.yt.orig.xml")
		if readErr != nil {
			return nil, readErr
		}

		return rss.Parse(data)
	}

	return rss.Fetch(*url) //"https://www.youtube.com/feeds/videos.xml?channel_id=UCWfRKs8owsEkERlwO1uwOFw"
}

func parseYtRss(ytFeed *rss.Feed) (p Podcast) {

	p.YtID = ytFeed.ID
	p.Lang = "ru-RU"
	p.Title = ytFeed.Title
	p.Link = ytFeed.Link
	p.AuthorName = ytFeed.Nickname
	p.Cached = time.Now()

	// Add categories
	catStrings := []string{"Society &amp; Culture/Personal Journals", "Technology/Tech News"}
	categories := []*Category{}
	for _, cat := range catStrings {
		tehCat := Category{}
		tehCat.Name = cat
		categories = append(categories, &tehCat)
	}
	p.Categories = categories

	// Add episodes
	episodes := []Episode{}
	for _, ytEp := range ytFeed.Items {

		ep := Episode{}

		ep.YtID = ytEp.ID
		ep.Title = ytEp.Title
		ep.Published = ytEp.Date
		ep.YtLink = ytEp.Link

		if len(ytEp.Desc) == 0 {
			ep.Description = "<No Shownotes>"
		} else {
			ep.Description = ytEp.Desc
		}

		author := Author{}
		author.Name = ytEp.Author

		ep.Author = author

		ep.Views = ytEp.Views

		ep.PodcastYtID = p.YtID
		episodes = append(episodes, ep)
	}
	p.Episodes = episodes

	return
}

func itcPodcastFrom(p *Podcast) podcast.Podcast {

	itcPod := podcast.New(p.Title, p.Link, p.Description, &p.FirstPublished, &p.Cached)

	itcPod.IAuthor = p.Title //AuthorName
	itcPod.Language = "ru-RU"
	itcPod.IExplicit = "false"

	for _, ep := range p.Episodes {
		itcItem := new(podcast.Item)

		itcItem.Title = ep.Title
		itcItem.PubDate = &ep.Published

		itcItem.Link = ep.YtLink
		itcItem.Description = ep.Description

		author := podcast.Author{}
		author.Name = p.AuthorName
		itcItem.Author = &author

		itcItem.Comments = strconv.Itoa(ep.Views) + " Views"

		_, addErr := itcPod.AddItem(*itcItem)
		if addErr != nil {
			log.Fatal(addErr)
		}
	}

	return itcPod
}

func writeItunesPodcastRssXML(itcPodcast podcast.Podcast) error {
	return ioutil.WriteFile("rsss/rss.itc.gen.xml", itcPodcast.Bytes(), 0644)
}
