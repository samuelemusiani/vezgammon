package config

import (
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Bgweb    Bgweb    `json:"bgweb"`
	Docker   bool     // true if is deployen with docker
	Swagger  bool     // if true expose swagger console
}

type Server struct {
	Bind   string `json:"bind"`
	Domain string `json:"domain"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type Bgweb struct {
	Url string `json:"domain"`
}

var conf = Config{
	Server{Bind: ":3001", Domain: "localhost:3001"},
	Database{User: "", Password: "", Address: ":5432"},
	Bgweb{Url: "localhost:3002"},
	false,
	false,
}

func Get() *Config {
	return &conf
}

func Set(c *Config) {
	conf = *c
}

func Parse(path string) error {
	buff, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(buff, &conf)
	if err != nil {
		return err
	}

	slog.With("conf", conf).Debug("")

	return nil
}
