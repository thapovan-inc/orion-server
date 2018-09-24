package util

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"reflect"
)

type ServerConfig struct {
	Logger          LoggerConfig         `toml:"log"`
	SpanValidation  SpanValidationConfig `toml:"span_validation"`
	PublisherConfig PublisherConfig      `toml:"publisher"`
	loaded          bool
}

type GeneralConfig struct {
}

type LoggerConfig struct {
	Level     string `toml:"level"`
	Format    string `toml:"format"`
	UseColors bool   `toml:"use_colors"`
}

type NatsPublisherConfig struct {
	URL       string `toml:"url"`
	ClientID  string `toml:"client_id"`
	ClusterID string `toml:"cluster_id"`
}

type SpanValidationConfig struct {
	AllowedDrift    int64  `toml:"allowed_drift_ms"`
	MaxMetaDataSize uint64 `toml:"max_metadata_size_kb"`
}

var loadedConfig ServerConfig

func GetConfig() ServerConfig {
	if !loadedConfig.loaded {
		panic("Config data not loaded")
	}
	return loadedConfig
}

func LoadConfig(tomlData string) {
	loadedConfig = ServerConfig{}
	_, err := toml.Decode(tomlData, &loadedConfig)
	if err != nil {
		panic(fmt.Errorf("Error when parsing toml data: %v", err))
	}
	if reflect.DeepEqual(ServerConfig{}, loadedConfig) {
		panic("Empty config data")
	} else {
		loadedConfig.loaded = true
	}
}

func LoadConfigFromFile(fileName string) {
	tomlData, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	LoadConfig(string(tomlData))
}
