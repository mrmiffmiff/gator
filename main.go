package main

import (
	"fmt"
	"os"

	"github.com/mrmiffmiff/gator-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading file: %w", err))
		os.Exit(1)
	}
	cfg.SetUser("robert")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading file: %w", err))
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)
}
