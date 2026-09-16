package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mrmiffmiff/gator-blog-aggregator/internal/database"
)

func scrapeFeeds(s *state) error {
	dbFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Trouble retrieving feed from database: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        dbFeed.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("Trouble marking feed as fetched: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return fmt.Errorf("Error with fetching feed data: %w", err)
	}
	fmt.Printf("RSS Feed: %s\n", rssFeed.Channel.Title)
	fmt.Println(rssFeed.Channel.Description)
	if len(rssFeed.Channel.Item) < 1 {
		fmt.Println("No items detected in RSS Feed")
		return nil
	}
	fmt.Println("Items are as follows:")
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Couldn't parse time between requests: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
