package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/R0Xps/gatorcli/internal/config"
	"github.com/R0Xps/gatorcli/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", conf.Db_url)

	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	st := state{
		db:     dbQueries,
		config: &conf,
	}
	cmds := commands{
		commands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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
	db     *database.Queries
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command 'login' expects 1 argument (username)")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])

	if err != nil {
		log.Fatal("User doesn't exist")
	}

	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Login successful")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command 'register' expects 1 argument (username)")
	}
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}

	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been created")
	fmt.Println(usr)
	return nil
}
