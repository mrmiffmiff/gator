package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	for _, item := range rssFeed.Channel.Item {
		var pubTime sql.NullTime
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				pubTime = sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				}
			} else {
				pubTime = sql.NullTime{
					Time:  t,
					Valid: true,
				}
			}
		} else {
			pubTime = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			Url:       item.Link,
			FeedID:    dbFeed.ID,
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubTime,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Duplicate insert
				continue
			}
			return fmt.Errorf("Error creating post: %w", err)
		}
		fmt.Printf("Added post with URL %s published at %s to database\n", post.Url, post.PublishedAt.Time.String())
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [post_limit]", cmd.Name)
	}
	if len(cmd.Args) < 1 {
		fmt.Println("No post limit entered, setting to 2.")
	}
	var limit string = "2"
	if len(cmd.Args) == 1 {
		limit = cmd.Args[0]
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("Trouble converting limit to int: %w", err)
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limitInt),
	})
	if err != nil {
		return fmt.Errorf("Trouble retrieving posts for user: %w", err)
	}
	if len(posts) < 1 {
		fmt.Println("No posts found")
		return nil
	}
	for _, post := range posts {
		var feedMessage string
		feed, err := s.db.GetFeedById(context.Background(), post.FeedID)
		if err != nil {
			feedMessage = "From an unidentified feed:"
		} else {
			feedMessage = fmt.Sprintf("From Feed %s:", feed.Name)
		}
		fmt.Println(feedMessage)
		fmt.Printf("Post with url %s\n", post.Url)
		if post.Title.Valid {
			fmt.Println(post.Title.String)
		}
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		if post.PublishedAt.Valid {
			fmt.Printf("Published at %v\n", post.PublishedAt.Time)
		}
	}
	return nil
}
