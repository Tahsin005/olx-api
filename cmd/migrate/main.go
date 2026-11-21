package main

import (
	"log"
	"os"

	"github.com/Tahsin005/olx-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if (len(os.Args) < 2) {
		log.Fatal("usage: migrate <up|down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
        "file://migrations",
        cfg.DatabaseUrl,
	)

	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	command := os.Args[1]
	switch command {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("invalid command. Use 'up' or 'down'")
	}
}