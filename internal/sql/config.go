package sql

import "time"

type Config struct {
	ConnString      string        `mapstructure:"conn_string"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
