package main

import (
	"context"
	"fmt"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Something wrong with feed fetching: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}
