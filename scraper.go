package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/AnhBigBrother/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concuren_nums int, duration time.Duration) {
	log.Printf("Scrapping on %v routines every %v duration", concuren_nums, duration)
	ticker := time.NewTicker(duration)
	for {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concuren_nums),
		)
		if err != nil {
			log.Println("Error on GetNextFeedToFetch:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, f := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, f)
		}
		wg.Wait()
		<-ticker.C
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error on parse date %v with err: %v:\n", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubTime,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique") {
				continue
			}
			log.Println("Failed to create post:", err)
		}
	}
	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
