package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config variables from toml file
type Config struct {
	Server    string
	Database  string
	SecretKey string
	LogPath   string
}

// Read and parse the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
