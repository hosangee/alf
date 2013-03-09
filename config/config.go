/*
	Package config implements the configuration parsing
*/
package config

import (
	"errors"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
)

//A Config is an implementation of the configuration parameters read
// from the config file
type Config struct {
	// DbUrl is the url used for connecting to the database. The syntax is
	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	DbUrl string

	// BindAddress is the address the http service will bind to. The syntax is
	// host:port
	BindAddress string

	file *yaml.File
}

// NewConfig returns a Config instance generated from the provided
// config file path.
func NewConfig(configPath string) (config *Config, err error) {
	parsedConfig, err := yaml.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config = &Config{}

	// Get database url
	value, err := parsedConfig.Get("db.url")
	if err != nil {
		return nil, errors.New("Invalid database configuration: " + err.Error())
	}
	config.DbUrl = strings.Replace(value, `"`, "", -1)

	// Get http url
	value, err = parsedConfig.Get("http.bind")
	if err != nil {
		return nil, errors.New("Invalid http configuration: " + err.Error())
	}
	config.BindAddress = strings.Replace(value, `"`, "", -1)

	return config, nil
}
