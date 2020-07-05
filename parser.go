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

func parseYtRss(ytFeed *rss.Feed) (pod Podcast) {

	pod.YtID = ytFeed.ID
	pod.Lang = "ru-RU"
	pod.Title = ytFeed.Title
	pod.Link = ytFeed.Link
	pod.AuthorName = ytFeed.Nickname

	catStrings := []string{"Society &amp; Culture/Personal Journals", "Technology/Tech News"}
	categories := []*Category{}
	for _, cat := range catStrings {
		tehCat := Category{}
		tehCat.Name = cat
		categories = append(categories, &tehCat)
	}
	pod.Categories = categories

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

		ep.PodcastYtID = pod.YtID
		episodes = append(episodes, ep)
	}

	pod.Episodes = episodes
	pod.Cached = time.Now()

	return 
}

func itcPodcastFrom(tehPod *Podcast) podcast.Podcast {

	p := podcast.New(tehPod.Title, tehPod.Link, tehPod.Description, &tehPod.FirstPublished, &tehPod.Cached)

	p.IAuthor = tehPod.Title //AuthorName
	p.Language = "ru-RU"
	p.IExplicit = "true"

	for _, ytEpisode := range tehPod.Episodes {
		itcItem := new(podcast.Item)

		itcItem.Title = ytEpisode.Title
		itcItem.PubDate = &ytEpisode.Published

		itcItem.Link = ytEpisode.YtLink
		itcItem.Description = ytEpisode.Description

		author := podcast.Author{}
		author.Name = tehPod.AuthorName
		itcItem.Author = &author

		itcItem.Comments = strconv.Itoa(ytEpisode.Views) + " Views"

		_, addErr := p.AddItem(*itcItem)
		if addErr != nil {
			log.Fatal(addErr)
		}
	}

	return p
}

func writeItunesPodcastRssXML(itcPodcast podcast.Podcast) error {
	return ioutil.WriteFile("rsss/rss.itc.gen.xml", itcPodcast.Bytes(), 0644)
}
