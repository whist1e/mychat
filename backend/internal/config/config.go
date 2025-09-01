package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type MysqlConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type LogConfig struct {
	LogPath    string `toml:"logPath"`
	Level      string `toml:"level"`
	Format     string `toml:"format"`
	MaxSize    int    `toml:"maxSize"`
	MaxBackups int    `toml:"maxBackups"`
	MaxAge     int    `toml:"maxAge"`
	Compress   bool   `toml:"compress"`
	Console    bool   `toml:"console"`
}

type Config struct {
	MainConfig  `toml:"mainConfig"`
	MysqlConfig `toml:"mysqlConfig"`
	RedisConfig `toml:"redisConfig"`
	LogConfig   `toml:"logConfig"`
}

var config *Config

func LoadConfig() error {
	// 使用相对路径，指向项目根目录下的configs/config.toml
	if _, err := toml.DecodeFile("configs/config.toml", config); err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = LoadConfig()
	}
	return config
}
