package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	callback, exist := c.registeredCommands[cmd.Name]
	if !exist {
		return fmt.Errorf("Command does not exist")
	}
	return callback(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
