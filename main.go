package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/mrmiffmiff/gator-blog-aggregator/internal/config"
	"github.com/mrmiffmiff/gator-blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	programState := &state{
		cfg: &cfg,
	}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	db, err := sql.Open("postgres", programState.cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	programState.db = dbQueries
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
