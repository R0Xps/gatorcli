package main

import (
	"fmt"
	"log"
	"os"

	"github.com/R0Xps/gatorcli/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	st := state{
		config: &conf,
	}
	cmds := commands{
		commands: map[string]func(*state, command) error{},
	}

	args := os.Args

	if len(args) < 2 {
		log.Fatal("Missing argument (command name)")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandHandler, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("command '%s' not found", cmd.name)
	}
	return commandHandler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
