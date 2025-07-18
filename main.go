package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const baseCfgPath = "myfeed/config.toml"

type Config struct {
	Resources []ResourceConfig `toml:"resources"`
}

type ResourceConfig struct {
	FeedURL string `toml:"feed_url"`
	ParserT string `toml:"parser"`
}

type ParserT = string

var (
	WebParserT      = ParserT("web")
	TelegramParserT = ParserT("telegram")
	TorrentParserT  = ParserT("torrent")
)

func main() {
	if os.Getenv("DEBUG") != "" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	var cfgPath string
	flag.StringVar(&cfgPath, "config", defaultCfg(), "path to a TOML config")
	flag.Parse()
	dat, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("failed to find config file at %s", cfgPath)
	}

	var conf Config
	_, err = toml.Decode(string(dat), &conf)
	if err != nil {
		log.Fatalf("failed to decode config at %s", cfgPath)
	}

	slog.Debug("initialization finished", "configured resources", len(conf.Resources))
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
