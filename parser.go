package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	rss "github.com/m4rr/yt-rss"
)

func parseYtRss(feed *rss.Feed) Podcast {

	thePod := Podcast{}

	thePod.YtID = feed.ID
	thePod.Lang = "ru-RU"
	thePod.Title = feed.Title
	thePod.Link = feed.Link
	thePod.AuthorName = feed.Nickname

	cats := []string{"Society &amp; Culture/Personal Journals", "Technology/Tech News"}
	cats3 := []*Category{}
	for _, cat := range cats {
		tehCat := Category{}
		tehCat.Name = cat
		cats3 = append(cats3, &tehCat)
	}
	thePod.Categories = cats3

	theEps := []Episode{}
	for _, ytEpisode := range feed.Items {
		theEp := Episode{}

		theEp.YtID = ytEpisode.ID
		theEp.Title = ytEpisode.Title
		theEp.Published = ytEpisode.Date
		theEp.YtLink = ytEpisode.Link

		if len(ytEpisode.Desc) == 0 {
			theEp.Description = "<No Shownotes>"
		} else {
			theEp.Description = ytEpisode.Desc
		}

		author := Author{}
		author.Name = ytEpisode.Author
		theEp.Author = author

		theEp.Views = ytEpisode.Views

		theEp.PodcastYtID = thePod.YtID
		theEps = append(theEps, theEp)
	}

	thePod.Episodes = theEps
	thePod.Cached = time.Now()

	return thePod
}

func itcPodcastFrom(thePod *Podcast) podcast.Podcast {

	p := podcast.New(thePod.Title, thePod.Link, thePod.Description, &thePod.FirstPublished, &thePod.Cached)

	p.IAuthor = thePod.Title //AuthorName
	p.Language = "ru-RU"
	p.IExplicit = "true"

	for _, ytEpisode := range thePod.Episodes {
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
