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

	"github.com/mmcdole/gofeed"
	_ "modernc.org/sqlite"

	"github.com/scipunch/myfeed/config"
	"github.com/scipunch/myfeed/parser"
	"github.com/scipunch/myfeed/parser/factory"
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
	var parserTypes []parser.Type
	for _, r := range conf.Resources {
		parserTypes = append(parserTypes, r.ParserT)
	}
	parsers, err := factory.Init(parserTypes)
	if err != nil {
		log.Fatalf("failed to initialize some parsers with %s", err)
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

	// Fetch configured RSS feeds
	var errs []error
	feeds := make([]*gofeed.Feed, len(conf.Resources))
	fp := gofeed.NewParser()
	for i, resource := range conf.Resources {
		feed, err := fp.ParseURL(resource.FeedURL)
		if err != nil {
			errs = append(errs, fmt.Errorf("'%s' parse failed with %w", resource.FeedURL, err))
			continue
		}
		feeds[i] = feed
	}
	if len(errs) > 0 {
		slog.Error("several feeds were not parsed", "feeds", errors.Join(errs...))
	}

	// Process new items
	errs = nil
	slog.Info("fetched feeds", "amount", len(feeds))
	for i, feed := range feeds {
		if feed == nil {
			slog.Debug("skipping failed to parse feed")
			continue
		}
		resource := conf.Resources[i]
		parser := parsers[resource.ParserT]
		for _, item := range feed.Items {
			data, err := parser.Parse(item.Link)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			slog.Info("feed item parsed", "length", len(data.String()))
		}
	}
	if len(errs) > 0 {
		slog.Error("failed to parse some pages", "errors", errors.Join(errs...).Error())
	}

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
