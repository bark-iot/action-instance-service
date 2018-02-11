package main

import (
	"github.com/urfave/cli"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"fmt"
	"os"
	"database/sql"
)

var (
	app        *cli.App
	db         *sql.DB
	mgr        *migrate.Migrate
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "migrate"
	app.Usage = "migrate action instances service tables"
	app.Author = "Serge Koba"
	app.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "run migrations up",
			Action: func(c *cli.Context) error {
				if err := up(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "down",
			Usage: "run migrations down",
			Action: func(c *cli.Context) error {
				if err := down(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func initMigrate() error {
	connectString := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@db:5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Printf("Cannot connect to database %v %v", connectString, err)
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "go_ai_schema_migrations"})
	if err != nil {
		fmt.Printf("Cannot init driver %v", err)
		db.Close()
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///go/src/app/migrate/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Printf("Cannot init migrate %v", err)
		db.Close()
		return err
	}
	mgr = m
	return nil
}

func up() error {
	if err := initMigrate(); err != nil {
		return err
	}
	if err := mgr.Up(); err != nil {
		fmt.Printf("Cannot migrate %v", err)
		return err
	}
	return nil
}

func down() error {
	if err := initMigrate(); err != nil {
		return err
	}
	if err := mgr.Down(); err != nil {
		fmt.Printf("Cannot migrate %v", err)
		return err
	}
	return nil
}