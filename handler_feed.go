package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrmiffmiff/gator-blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current user: %w", err)
	}
	currentUserId := currentUser.ID
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUserId,
	})
	if err != nil {
		return fmt.Errorf("Problem creating new feed: %w", err)
	}
	fmt.Println("New feed record created successfully")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds: %w", err)
	}
	fmt.Println("Feeds List:")
	for i, feed := range feeds {
		fmt.Printf("Feed %d\n", i+1)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't retrieve user: %w", err)
		}
		userName := user.Name
		fmt.Printf("* User Name:     %s\n", userName)
		fmt.Println("========")
	}
	return nil
}
