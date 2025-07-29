package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)

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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Database reset successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Print(user.Name)
		if user.Name == s.config.Current_user_name {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feed := RSSFeed{}
	if err = xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil

}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("command 'addfeed' expects 2 arguments (feedName, feedURL)")
	}

	user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		return err
	}
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command 'follow' expects 1 argument (url)")
	}

	user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(ff.FeedName, ff.UserName)
	return nil
}
