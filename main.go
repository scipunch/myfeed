package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"

	_ "modernc.org/sqlite"

	"github.com/scipunch/myfeed/config"
)

//go:embed schema.sql
var ddl string

func main() {
	if os.Getenv("DEBUG") != "" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	var cfgPath string
	flag.StringVar(&cfgPath, "config", config.DefaultPath(), "path to a TOML config")
	flag.Parse()

	// Read config and create if default is missing
	conf, err := config.Read(cfgPath)
	if errors.Is(err, os.ErrNotExist) && cfgPath == config.DefaultPath() {
		if err := config.Write(cfgPath, conf); err != nil {
			log.Fatalf("failed to write default config with %s", err)
		}
	} else if err != nil {
		log.Fatalf("failed to read config with %s", err)
	}

	// Connect to database & initialize schema
	dbBasePath := path.Dir(conf.DatabasePath)
	err = os.MkdirAll(dbBasePath, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create base shared directory at '%s' with %s", dbBasePath, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_, err = initDB(ctx, conf.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database schema with %v", err)
	}

	// TODO: Fetch configured RSS feeds
	// TODO: Process new items
	// TODO: Generate PDF report
}

func initDB(ctx context.Context, source string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", source)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at '%s' with %w", source, err)
	}
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, fmt.Errorf("failed to execute DDL with %w", err)
	}
	return db, nil
}
