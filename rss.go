package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func urlToFeed(url string) (RSSFeed, error) {
	rssFeed := RSSFeed{}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return rssFeed, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return rssFeed, err
	}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return rssFeed, err
	}
	return rssFeed, nil
}