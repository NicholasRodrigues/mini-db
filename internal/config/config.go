package config

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Server struct {
		Port        string `mapstructure:"port"`
		TLS         bool   `mapstructure:"tls"`
		TLSCertFile string `mapstructure:"tls_cert_file"`
		TLSKeyFile  string `mapstructure:"tls_key_file"`
	} `mapstructure:"server"`
	Storage struct {
		FilePath string `mapstructure:"file_path"`
	} `mapstructure:"storage"`
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	Security struct {
		AuthEnabled bool   `mapstructure:"auth_enabled"`
		AuthToken   string `mapstructure:"auth_token"`
	} `mapstructure:"security"`
}

// Cfg is the global configuration instance
var Cfg Config

// LoadConfig reads configuration from file or environment variables
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	if Cfg.Logging.Level == "" {
		Cfg.Logging.Level = "info"
	}

	logLevel, err := logrus.ParseLevel(strings.ToLower(Cfg.Logging.Level))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(logLevel)
}
