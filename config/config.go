package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type JWTConfig struct {
	Secret          string        `yaml:"key"`
	Issuer          string        `yaml:"issuer"`
	Audience        string        `yaml:"audience"`
	ExpirationHours int           `yaml:"expiration_hours"`
	RefreshDays     int           `yaml:"refresh_days"`
	Expiration      time.Duration `yaml:"-"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
type log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

func LoadFromFile(configPath string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}
	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 解析 YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	config.JWT.Expiration = time.Duration(config.JWT.ExpirationHours) * time.Hour
	return config, nil
}
