package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type contextKey uint

const (
	RequestID contextKey = iota
	FieldUUID
)

type Config struct {
	Env        string `yaml:"env" env-default:"production"`
	Storage    string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Type        string        `yaml:"type" env-default:"net/http"`
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.Env),
		slog.String("storage", c.Storage),
		slog.Group("http_server",
			slog.String("type", c.Type),
			slog.String("address", c.Address),
			slog.Duration("timeout", c.Timeout),
			slog.Duration("idle_timeout", c.IdleTimeout),
			slog.String("user", c.User),
			slog.String("password", "***"),
		),
	)
}

func NewConfig() *Config {
	cfgPath := fetchConfigPath()
	if cfgPath == "" {
		panic("CONFIG_PATH is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
