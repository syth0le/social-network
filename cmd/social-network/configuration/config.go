package configuration

import "social-network/internal/utils"

type Config struct {
	Logger  *LoggerConfig  `yaml:"logger"`
	Server  *ServerConfig  `yaml:"server"`
	Storage *StorageConfig `yaml:"storage"`
}

func (c *Config) Validate() error {
	return nil // todo
}

type LoggerConfig struct {
	Level      utils.Level `yaml:"level"`
	Encoding   string      `yaml:"encoding"`
	Path       string      `yaml:"path"`
	App        string      `yaml:"app"`
	Enviroment string      `yaml:"enviroment"`
}

func (c *LoggerConfig) Validate() error {
	return nil // todo
}

type ServerConfig struct {
	Enable   bool   `yaml:"enable_mock"`
	Endpoint string `yaml:"endpoint`
	Port     string `yaml:"port"`
}

func (c *ServerConfig) Validate() error {
	return nil // todo
}

type StorageConfig struct {
	EnableMock         bool   `yaml:"enable_mock"`
	Host               string `yaml:"host" env:"HOST"`
	Port               int    `yaml:"port" env:"PORT"`
	Database           string `yaml:"database" env:"DATABASE"`
	Username           string `yaml:"username" env:"USER"`
	Password           string `yaml:"password" env:"PASSWORD"`
	SSLMode            string `yaml:"ssl_mode"`
	ConnectionAttempts int    `yaml:"connection_attempts"`
}

func (c *StorageConfig) Validate() error {
	return nil // todo
}
