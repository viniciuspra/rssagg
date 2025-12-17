package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/viniciuspra/rssagg/internal/db"
)

func startScraping(ctx context.Context, db *db.Queries, concurrency int, timeBtwReq time.Duration) {
	log.Printf("starting scraping on %v goroutines every %v", concurrency, timeBtwReq)

	ticker := time.NewTicker(timeBtwReq)
	defer ticker.Stop()
	for {
		select {
			case <-ctx.Done():
				log.Println(ctx.Err())
				return
			case <-ticker.C:
				runScrapeCycle(ctx, db, concurrency)
		}
	}
}

func runScrapeCycle(ctx context.Context, db *db.Queries, concurrency int) {
	wg := &sync.WaitGroup{}
	feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
	if err != nil {
		log.Println("error fetching next feeds to fetch:", err)
		return
	}
	for _, feed := range feeds {
		wg.Add(1)
		go scrapeFeed(ctx, wg, db, feed)
	}
	wg.Wait()
}

func scrapeFeed(ctx context.Context, wg *sync.WaitGroup, db *db.Queries, feed db.Feed) {
	defer wg.Done()

	err := db.MarkFeedAsFetched(ctx, feed.ID)
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
