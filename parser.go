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

func parseYtRss(feed *rss.Feed) Podcast {

	tehPod := Podcast{}

	tehPod.YtID = feed.ID
	tehPod.Lang = "ru-RU"
	tehPod.Title = feed.Title
	tehPod.Link = feed.Link
	tehPod.AuthorName = feed.Nickname

	cats := []string{"Society &amp; Culture/Personal Journals", "Technology/Tech News"}
	cats3 := []*Category{}
	for _, cat := range cats {
		tehCat := Category{}
		tehCat.Name = cat
		cats3 = append(cats3, &tehCat)
	}
	tehPod.Categories = cats3

	tehEps := []Episode{}
	for _, ytEpisode := range feed.Items {
		tehEp := Episode{}

		tehEp.YtID = ytEpisode.ID
		tehEp.Title = ytEpisode.Title
		tehEp.Published = ytEpisode.Date
		tehEp.YtLink = ytEpisode.Link

		if len(ytEpisode.Desc) == 0 {
			tehEp.Description = "<No Shownotes>"
		} else {
			tehEp.Description = ytEpisode.Desc
		}

		author := Author{}
		author.Name = ytEpisode.Author
		tehEp.Author = author

		tehEp.Views = ytEpisode.Views

		tehEp.PodcastYtID = tehPod.YtID
		tehEps = append(tehEps, tehEp)
	}

	tehPod.Episodes = tehEps
	tehPod.Cached = time.Now()

	return tehPod
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
