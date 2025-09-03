package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

type GRPCConfig struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env-default:"50051"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	DSN string        `yaml:"dsn"`
	TTL time.Duration `yaml:"ttl" env-default:"1h"`
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
