package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	DB         Database   `yaml:"db"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type Database struct {
	DBName   string `yaml:"dbname" env:"DATABASE_NAME"`
	User     string `yaml:"user" env:"DATABASE_USER"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DATABASE_PORT"`
	SSLMode  string `yaml:"sslmode" env:"DATABASE_SSLMODE"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH environment variable is not set.")
	}

	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Config path does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config.")
	}

	return &cfg
}
