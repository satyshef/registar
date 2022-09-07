// Package config конфигурация GoMon
package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIKey  string `toml:"api_key"`
	ID      string `toml:"id"`
	Country string `toml:"country"`
	Host    string `toml:"host"`
}

func LoadConfig(configPath string) (*Config, error) {
	c := New()
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при загрузке конфигурации : %s", err)

	}

	return c, nil
}

// New инициализация конфигурации приложения
func New() *Config {
	return &Config{}
}
