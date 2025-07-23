package config

import (
	"os"
	"path"
)

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

func Default() Config {
	var dbBase = path.Join(os.Getenv("HOME"), ".local/share/myfeed")
	return Config{
		DatabasePath: path.Join(dbBase, "data.db"),
	}
}
