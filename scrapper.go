package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/viniciuspra/rssagg/internal/db"
)

func startScraping(ctx context.Context, dbQ *db.Queries, concurrency int, timeBtwReq time.Duration) {
	log.Printf("starting scraping on %v goroutines every %v", concurrency, timeBtwReq)

	ticker := time.NewTicker(timeBtwReq)
	defer ticker.Stop()
	runScrapeCycle(ctx, dbQ, concurrency)
	for {
		select {
			case <-ctx.Done():
				log.Println(ctx.Err())
				return
			case <-ticker.C:
				runScrapeCycle(ctx, dbQ, concurrency)
		}
	}
}

func runScrapeCycle(ctx context.Context, dbQ *db.Queries, concurrency int) {
	wg := &sync.WaitGroup{}
	feeds, err := dbQ.GetNextFeedsToFetch(ctx, int32(concurrency))
	if err != nil {
		log.Println("error fetching next feeds to fetch:", err)
		return
	}
	for _, feed := range feeds {
		wg.Add(1)
		go scrapeFeed(ctx, wg, dbQ, feed)
	}
	wg.Wait()
}

func scrapeFeed(ctx context.Context, wg *sync.WaitGroup, dbQ *db.Queries, feed db.Feed) {
	defer wg.Done()

	err := dbQ.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching next feeds to fetch:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		select {
	    case <-ctx.Done():
	        return
	    default:
    }
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Println("error parsing date:", err)
			continue
		}

		_, err = dbQ.CreatePost(ctx, db.CreatePostParams{
			ID: uuid.New(),
			FeedID: feed.ID,
			Title: item.Title,
			Description: description,
			PublishedAt: pubDate,
			Url: item.Link,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("error creating post %v with err: %v", item.Title, err)
		}
	}

	log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func parsePubDate(pubDate string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
	}

	var parsedDate time.Time
	var err error

	for _, layout := range layouts {
		parsedDate, err = time.Parse(layout, pubDate)
		if err == nil {
			return parsedDate, nil
		}
	}

	return time.Time{}, err
}
