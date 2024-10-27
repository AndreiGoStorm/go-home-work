package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger     LoggerConf
	HTTPServer HTTPServer
	Storage    Storage
	Database   Database
}

type LoggerConf struct {
	Level string
}

type HTTPServer struct {
	Host string
	Port int
}

type Storage struct {
	Driver string
}

type Database struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
}

func New(path string) *Config {
	config := &Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("create config error: %w", err))
	}
	return config
}
