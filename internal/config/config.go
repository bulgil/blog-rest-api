package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	PGStorage  PGStorage  `yaml:"pg_storage"`
}

type PGStorage struct {
	Address  string `yaml:"address" env-default:"127.0.0.1"`
	Port     string `yaml:"port" env-default:"10000"`
	Login    string `yaml:"login" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"127.0.0.1"`
	Port    string `yaml:"port" env-default:"10000"`
}

func NewConfig() *Config {
	configPath := flag.String("config", "./config/config.yml", "path to config file")
	flag.Parse()

	var cfg Config

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}

	log.Println("config initialized")
	return &cfg
}
