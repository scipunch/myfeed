package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const baseCfgPath = "myfeed/config.toml"

type Config struct {
	Resources    []ResourceConfig `toml:"resources"`
	DatabasePath string           `toml:"database_path"`
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

func Read(path string) (Config, error) {
	conf := Default()
	if path != DefaultPath() {
		dat, err := os.ReadFile(path)
		if err != nil {
			return conf, fmt.Errorf("failed to find config file at %s with %w", path, err)
		}
		_, err = toml.Decode(string(dat), &conf)
		if err != nil {
			return conf, fmt.Errorf("failed to decode config at %s with %w", path, err)
		}
		slog.Debug("initialization finished", "configured resources", len(conf.Resources))
	}
	return conf, nil
}

func Default() Config {
	var dbBase = path.Join(os.Getenv("HOME"), ".local/share/myfeed")
	return Config{
		DatabasePath: path.Join(dbBase, "data.db"),
	}
}

func DefaultPath() string {
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
