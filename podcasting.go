package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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

	splitID := strings.Split(ytFeed.ID, ":")
	p.Nickname = splitID[len(splitID)-1]

	// Add categories
	// catStrings := []string{"Society & Culture/Personal Journals", "Technology/Tech News"}
	// categories := []Category{}
	// for _, cat := range catStrings {
	// 	tehCat := Category{}
	// 	tehCat.Name = cat
	// 	categories = append(categories, &tehCat)
	// }
	// p.Categories = categories

	// Add episodes
	episodes := []Episode{}
	for _, ytEp := range ytFeed.Items {

		ep := Episode{}

		ep.YtID = ytEp.ID
		// ep.ChannelID = ytFeed.

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

func fillCategories() []Category {
	categories := []Category{}
	lastSuper := ""
	for _, cat := range categoriesList() {
		tehCat := Category{}

		if strings.HasPrefix(cat, "-") {
			lastSuper = strings.TrimPrefix(cat, "-")
			tehCat.Name = lastSuper
		} else {
			tehCat.Name = lastSuper + "/" + cat
		}

		categories = append(categories, tehCat)
	}
	return categories
}

func categoriesList() []string {
	return []string{
		"-Arts",
		"Design",
		"Fashion & Beauty",
		"Food",
		"Literature",
		"Performing Arts",
		"Visual Arts",
		"-Business",
		"Business News",
		"Careers",
		"Investing",
		"Management & Marketing",
		"Shopping",
		"-Comedy",
		"-Education",
		"Education Technology",
		"Higher Education",
		"K-12",
		"Language Courses",
		"Training",
		"-Games & Hobbies",
		"Automotive",
		"Aviation",
		"Hobbies",
		"Other Games",
		"Video Games",
		"-Government & Organizations",
		"Local",
		"National",
		"Non-Profit",
		"Regional",
		"-Health",
		"Alternative Health",
		"Fitness & Nutrition",
		"Self-Help",
		"Sexuality",
		"-Kids & Family",
		"-Music",
		"-News & Politics",
		"-Religion & Spirituality",
		"Buddhism",
		"Christianity",
		"Hinduism",
		"Islam",
		"Judaism",
		"Other",
		"Spirituality",
		"-Science & Medicine",
		"Medicine",
		"Natural Sciences",
		"Social Sciences",
		"-Society & Culture",
		"History",
		"Personal Journals",
		"Philosophy",
		"Places & Travel",
		"-Sports & Recreation",
		"Amateur",
		"College & High School",
		"Outdoor",
		"Professional",
		"-Technology",
		"Gadgets",
		"Podcasting",
		"Software How-To",
		"Tech News",
		"-TV & Film",
	}
}
