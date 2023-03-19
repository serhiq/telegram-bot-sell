package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

const TempPatch = "./tmp/"
const PreviewCachePatch = "./imageCache/"

type Config struct {
	Telegram      Telegram      `yaml:"telegram"`
	RestaurantAPI RestaurantAPI `yaml:"restaurant_api"`
	DBConfig      DBConfig      `yaml:"database"`
}

type Telegram struct {
	Token string `yaml:"token" envconfig:"TELEGRAM_TOKEN" validate:"required"`
}

type RestaurantAPI struct {
	BaseURL string `yaml:"base_url" envconfig:"RESTAURANT_API_BASE_URL,omitempty"`
	Auth    string `yaml:"auth" envconfig:"RESTAURANT_API_AUTH,omitempty"`
	Store   string `yaml:"store" envconfig:"RESTAURANT_API_STORE,omitempty"`
}

type DBConfig struct {
	Host         string `yaml:"host" envconfig:"DB_HOST"`
	Port         int    `yaml:"port" envconfig:"DB_PORT"`
	DatabaseName string `yaml:"database_name" envconfig:"DB_DATABASE_NAME"`
	Username     string `yaml:"username" envconfig:"DB_USERNAME"`
	Password     string `yaml:"password" envconfig:"DB_PASSWORD"`
}

func New() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/config.yaml"
	}

	config := &Config{}

	if err := fromYaml(path, config); err != nil {
		fmt.Printf("couldn'n load config from %s: %s\r\n", path, err.Error())
	}

	if err := fromEnv(config); err != nil {
		fmt.Printf("couldn'n load config from env: %s\r\n", err.Error())
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(TempPatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating tmp path %s (%s)", filepath.Dir(TempPatch), err)
	}

	if err := os.MkdirAll(filepath.Dir(PreviewCachePatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating cache path %s (%s)", filepath.Dir(TempPatch), err)
	}

	return config, nil
}

func fromYaml(path string, config *Config) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func fromEnv(config *Config) error {
	return envconfig.Process("", config)
}

func validate(cfg *Config) error {
	if cfg.Telegram.Token == "" {
		return fmt.Errorf("config: %s is not set", "TELEGRAM_TOKEN")
	}

	if cfg.RestaurantAPI.BaseURL == "" {
		return fmt.Errorf("config: %s is not set", "RESTAURANT_API_BASE_URL")
	}

	if cfg.RestaurantAPI.Auth == "" {
		return fmt.Errorf("config: %s is not set", "RESTAURANT_API_AUTH")
	}

	if cfg.DBConfig.DatabaseName == "" {
		return fmt.Errorf("config: %s is not set", "DB_DATABASE_NAME")
	}

	if cfg.DBConfig.Username == "" {
		return fmt.Errorf("config: %s is not set", "DB_USERNAME")
	}

	return nil
}
