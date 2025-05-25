package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"service/internal/pkg/logs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var op = "cmd.migrator"

func main() {
	file, w, err := logs.Setup(context.Background(), "./files/logs/migrator/")
	if err != nil {
		return
	}
	defer file.Close()
	defer w.Flush()

	var storageURL, action, migrationsPath, migrationsTable string

	flag.StringVar(&storageURL, "storage-URL", "", "")
	flag.StringVar(&action, "action", "", "up/down/revert")
	flag.StringVar(&migrationsPath, "migrations-path", "", "")
	flag.StringVar(&migrationsTable, "migrations-table", "", "")
	flag.Parse()

	migration, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("postgres://%s?x-migrations-table=%s&sslmode=enable", storageURL, migrationsTable))
	if err != nil {
		logs.Error(
			context.Background(),
			err.Error(),
			op,
		)

		return
	}

	switch action {
	case "up":
		if err = migration.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logs.Info(
					context.Background(),
					"no new migrations",
					op,
				)
			} else {
				logs.Error(
					context.Background(),
					err.Error(),
					op,
				)
				return
			}
		}
	case "down":
		if err := migration.Down(); err != nil {
			logs.Error(
				context.Background(),
				err.Error(),
				op,
			)

			return
		}
	}

	logs.Info(
		context.Background(),
		"migrations applied",
		op,
	)
}
