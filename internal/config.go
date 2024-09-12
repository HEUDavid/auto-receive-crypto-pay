package internal

import (
	"github.com/BurntSushi/toml"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"log"
	"path/filepath"
)

type Config struct {
	Global       Global                     `toml:"global"`
	AdminAddress map[string][]AddressConfig `toml:"adminAddress"`
}

type Global struct {
	Mode     string `toml:"mode"`
	LogPath  string `toml:"logPath"`
	Addr     string `toml:"addr"`
	HostRoot string `toml:"hostRoot"`
	Auth     string `toml:"auth"`
}

type AddressConfig struct {
	Address string `toml:"address"`
	URL     string `toml:"url"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		if _, err := toml.DecodeFile(filepath.Join(util.FindProjectRoot(), "conf", "conf.toml"), &config); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}
	return config
}
