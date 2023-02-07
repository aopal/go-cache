package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port    string   `yaml:"port"`
	Origins []string `yaml:"origins"`

	Cache CacheConfig `yaml:"cache"`
	// FetchWorkers int      `yaml:"fetch_worker_count"`
}

type CacheConfig struct {
	QueryStringSort bool `yaml:"sort_query_string"`
}

func New(configPath string) (*Config, error) {
	var config Config

	yamlBytes, err := ioutil.ReadFile(configPath)
	if err != nil || len(yamlBytes) == 0 {
		return nil, fmt.Errorf("could not load config file '%+v': %+v", configPath, err)
	}

	yamlBytes = []byte(os.ExpandEnv(string(yamlBytes)))

	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%+v': %+v", configPath, err)
	}

	return &config, nil
}
