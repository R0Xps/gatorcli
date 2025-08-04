# Gator
Gator is a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View titles of the aggregated posts in the terminal, with a link to the full post

## Requirements
- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/) database
- [Goose](https://github.com/pressly/goose)

## Setting up requirements
### Installing Go
There are 2 ways of installing go:
- Using [Webi](https://webinstall.dev/golang/), this is the easiest way to do it in my opinion

- Or you could go to the [official Go installation page](https://go.dev/doc/install) and follow the instructions there for your platform

### Installing PostgreSQL
To install PostgreSQL use the following commands:
```bash
# on MacOS
brew install postgresql@15

# on Linux/WSL (Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Then to verify the installation:
```bash
psql --version
```

And finally, if you're on Linux, you need to update the postgres user password with the following command:
```bash
sudo passwd postgres
```
Use a password that you can easily remember, for example, I set it to 'postgres'.

### Running a Postgres server and creating a database
Start the Postgres server in the background:
```bash
# Mac 
brew services start postgresql@15

# Linux
sudo service postgresql start
```

Then connect to the Postgres server:
```bash
# Mac
psql postgres

# Linux
sudo -u postgres psql
```

Now we can create our database:
```PostgreSQL
CREATE DATABASE gator;
```

Then to make sure the database was created:
```PostgreSQL
\l
```
And make sure you see your gator table in the list.

If you are on Linux, make sure you change the postgres user password from here as well (this will be the password used in your connection string later):
```PostgreSQL
ALTER USER postgres PASSWORD 'postgres';
```

And finally after the database is created, you can disconnect from the Postgres server:
```PostgreSQL
exit
```

### Grabbing the database connection string
You will need this string in later steps to its better to take note of it now before you forget your database credentials.

The string is formatted as follows: `protocol://username:password@host:port/database`

So the full connection string should look like this:

 If you're on Mac (no password, your username): `postgres://your_username:@localhost:5432/gator`

And if you're on Linux (username is postgres, password is whatever you set it to when you installed Postgres): `postgres://postgres:postgres@localhost:5432/gator`

### Installing Goose
For more information about Goose, you can check out the README.md on their [GitHub repo](https://github.com/pressly/goose).

But for now we only need to install it and that's as simple as running this command in your terminal:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure Goose was installed correctly with the following command:
```bash
goose -version
```

### Preparing the database
First of all, you need to clone this GitHub repo, and to do that, create a new directory wherever you want it to be on your system (you only need the repo to set up the database, so you can delete the directory when you're done with this step)

```bash
# Create a new directory called 'gator'
mkdir gator

# Change the working directory to our new 'gator' directory
cd gator
```

Now clone the repo to your current directory:
```bash
git clone https://github.com/R0Xps/gatorcli .
```

Then move into the sql/schema directory:
```bash
cd sql/schema
```

Here comes the first time to use your connection string that you got after creating the database, the final step of preparing the database, running the goose migrations:
```bash
goose postgres <connection_string> up 
```

After this, everything should be ready to install and start using Gator!

## Installing Gator

To install the Gator CLI tool, you can simply run:
```bash
go install github.com/R0Xps/gatorcli/cmd/gator@latest
```

And before you can start using the tool, you have to create a config file in your home directory, it should be named `.gatorconfig.json` with the following contents (replace <connection_string> with your connection string that you got after creating the database)

```json
{
  "db_url": "<connection_string>?sslmode=disable"
}
```

## Running Gator
After installing Gator and creating a config file with the correct contents, you can use the tool by running the commands as shown in the next section.

Some commands require you to be registered and logged in, so it is recommended that you register before doing anything else.
You can register using the command `gator register <username>`

## Commands
To execute commands, you type them after the `gator` keyword in your terminal, following the template `gator command <arguments>`

The following is a list of all currently available commands and their usage:
- **register**: used to create new users. Usage: `register <username>`
- **login**: used to switch to an existing user. Usage `login <username>`
- **reset**: clear the database. Usage `reset`
- **users**: list all registered users indicating which is currently active. Usage `users`
- **agg**: start the aggregator, runs an infinite loop that grabs the posts from feeds stored in the database, with time between iterations given to the command. Usage `agg <time_between_requests>`
- **addfeed**: add a new feed to the database. Usage `addfeed <feed_name> <feed_url>`
- **feeds**: list all feeds in the database. Usage `feeds` 
- **follow**: follow a feed that has been added to the database by another user. Usage `follow <feed_url>`
- **following**: list all feeds followed by the currently active user. Usage `following`
- **unfollow**: unfollows a feed that you're following. Usage `follow <feed_url>`
- **browse**: list up to `limit` posts gathered by the aggregator from the `agg` command, defaults to 2 if not given. Usage `browse [limit]`