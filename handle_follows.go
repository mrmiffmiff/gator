package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrmiffmiff/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
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
	fmt.Printf("%s now follows %s\n", newFollow.UserName, newFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	userId := user.ID
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), userId)
	if err != nil {
		return fmt.Errorf("Error retrieving follows for user: %w", err)
	}
	if len(follows) < 1 {
		fmt.Printf("No feed follows for user %s.\n", user.Name)
		return nil
	}
	fmt.Printf("User %s follows the following feeds:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("  * %s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't retrieve feed: %w", err)
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing: %w", err)
	}
	fmt.Printf("%s unfollowed successfully for user %s!\n", feed.Name, user.Name)
	return nil
}
