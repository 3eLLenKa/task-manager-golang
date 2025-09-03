package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env"`
	Serv Server `yaml:"server"`
	GRPC GRPC   `yaml:"grpc"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Host            string        `yaml:"host"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type GRPC struct {
	DBService DBService `yaml:"db_service"`
}

type DBService struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Retry   Retry         `yaml:"retry"`
}

type Retry struct {
	Attempts   int           `yaml:"attempts"`
	Backoff    time.Duration `yaml:"backoff"`
	MaxBackoff time.Duration `yaml:"max_backoff"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic("config file is not exists")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("error to the read config file")
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
