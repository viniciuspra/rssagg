package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/viniciuspra/rssagg/internal/db"
)

func startScraping(db *db.Queries, concurrency int, timeBtwReq time.Duration) {
	log.Printf("starting scraping on %v goroutines every %v", concurrency, timeBtwReq)

	ticker := time.NewTicker(timeBtwReq)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching next feeds to fetch:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *db.Queries, feed db.Feed) {
	defer wg.Done()

	err := db.MarkFeedAsFetched(context.Background(), feed.ID)
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
		log.Printf("found item %s on feed (%s)\n", item.Title, feed.Name)
	}

	log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
