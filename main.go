package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/BurntSushi/toml"
	_ "modernc.org/sqlite"

	"github.com/scipunch/myfeed/config"
)

const baseCfgPath = "myfeed/config.toml"

//go:embed schema.sql
var ddl string

func main() {
	if os.Getenv("DEBUG") != "" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	var cfgPath string
	flag.StringVar(&cfgPath, "config", defaultCfg(), "path to a TOML config")
	flag.Parse()

	conf := config.Default()
	if cfgPath != defaultCfg() {
		dat, err := os.ReadFile(cfgPath)
		if err != nil {
			log.Fatalf("failed to find config file at %s", cfgPath)
		}

		_, err = toml.Decode(string(dat), &conf)
		if err != nil {
			log.Fatalf("failed to decode config at %s", cfgPath)
		}
		slog.Debug("initialization finished", "configured resources", len(conf.Resources))
	}

	dbBasePath := path.Dir(conf.DatabasePath)
	err := os.MkdirAll(dbBasePath, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create base shared directory at '%s' with %s", dbBasePath, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_, err = initDB(ctx, conf.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database schema with %v", err)
	}
}

func defaultCfg() string {
	var xdgHome = os.Getenv("XDG_CONFIG_HOME")
	if xdgHome != "" {
		return path.Join(xdgHome, baseCfgPath)
	}

	var home = os.Getenv("HOME")
	if home != "" {
		return path.Join(xdgHome, ".config", baseCfgPath)
	}

	panic("unclear where to search for the config fie")
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
