package app

import (
	"fmt"
	"go_service_template/internal/sql"

	"github.com/spf13/viper"
)

// Config - config of whole application
type Config struct {
	Database *sql.Config `mapstructure:"database"`
}

// ReadConfigFromFile - reads config from file
func ReadConfigFromFile(alterPath string) (*Config, error) {
	var config = &Config{}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if alterPath != "" {
		viper.AddConfigPath(alterPath)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config from file")
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return config, nil
}
