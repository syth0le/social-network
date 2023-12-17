package configuration

import (
	"social-network/internal/utils"
	"time"
)

type Config struct {
	Logger       LoggerConfig      `yaml:"logger"`
	Application  ApplicationConfig `yaml:"application"`
	PublicServer ServerConfig      `yaml:"public_server"`
	AdminServer  ServerConfig      `yaml:"admin_server"`
	Storage      StorageConfig     `yaml:"storage"`
}

func (c *Config) Validate() error {
	return nil // todo
}

type LoggerConfig struct {
	Level       utils.Level       `yaml:"level"`
	Encoding    string            `yaml:"encoding"`
	Path        string            `yaml:"path"`
	Environment utils.Environment `yaml:"environment"`
}

func (c *LoggerConfig) Validate() error {
	return nil // todo
}

type ApplicationConfig struct {
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout"`
	App                     string        `yaml:"app"`
	SaltValue               string        `yaml:"salt_value"`
}

func (c *ApplicationConfig) Validate() error {
	return nil // todo
}

type ServerConfig struct {
	Enable       bool   `yaml:"enable"`
	Endpoint     string `yaml:"endpoint"`
	Port         int    `yaml:"port" env:"PORT"`
	JwtTokenSalt string `yaml:"jwt_token_salt" env:"JWT_TOKEN_SALT"`
}

func (c *ServerConfig) Validate() error {
	return nil // todo
}

type StorageConfig struct {
	EnableMock         bool   `yaml:"enable_mock"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Database           string `yaml:"database"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password" env:"DB_PASSWORD"`
	SSLMode            string `yaml:"ssl_mode"`
	ConnectionAttempts int    `yaml:"connection_attempts"`
}

func (c *StorageConfig) Validate() error {
	return nil // todo
}
