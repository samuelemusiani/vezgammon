package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Server Server `json:"server"`
}

type Server struct {
	Bind string `json:"bind"`
}

var conf = Config{
	Server{Bind: ":3001"},
}

func getConf() *Config {
	return &conf
}

func parseConf(path string) error {
	buff, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(buff, &conf)
	if err != nil {
		return err
	}

	return nil
}
