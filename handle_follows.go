package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrmiffmiff/gator-blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	userName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current user: %w", err)
	}
	userId := user.ID
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed: %w", err)
	}
	feedId := feed.ID
	newFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		FeedID:    feedId,
	})
	if err != nil {
		return fmt.Errorf("Error creating new feed follow: %w", err)
	}
	fmt.Printf("%s now follows %s", newFollow.UserName, newFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	userName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve current user: %w", err)
	}
	userId := user.ID
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), userId)
	if err != nil {
		return fmt.Errorf("Error retrieving follows for user: %w", err)
	}
	if len(follows) < 1 {
		fmt.Printf("No feed follows for user %s.\n", userName)
		return nil
	}
	fmt.Printf("User %s follows the following feeds:\n", userName)
	for _, follow := range follows {
		fmt.Printf("  * %s\n", follow.FeedName)
	}
	return nil
}
