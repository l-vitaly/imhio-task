package config

import (
	"fmt"

	"github.com/l-vitaly/goenv"
)

// env name constants
const (
	HTTPAddrEnvName  = "CFGM_HTTP_ADDR"
	DBConnStrEnvName = "CFGM_DB_CONN_STR"
)

// Config service configuration
type Config struct {
	HTTPAddr string
	DB       struct {
		URL string
	}
}

// Get get env config vars
func Get() (*Config, error) {
	cfg := &Config{}
	goenv.StringVar(&cfg.HTTPAddr, HTTPAddrEnvName, ":9000")
	goenv.StringVar(&cfg.DB.URL, DBConnStrEnvName, "")
	goenv.Parse()

	if cfg.HTTPAddr == "" {
		return nil, fmt.Errorf("could not set %s", HTTPAddrEnvName)
	}
	if cfg.DB.URL == "" {
		return nil, fmt.Errorf("could not set %s", DBConnStrEnvName)
	}
	return cfg, nil
}
