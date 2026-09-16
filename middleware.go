package main

import (
	"context"
	"fmt"

	"github.com/mrmiffmiff/gator-blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		userName := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), userName)
		if err != nil {
			return fmt.Errorf("Couldn't retrieve current user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
