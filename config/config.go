package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
}

func LoadAndRead[T any](file string) (*T, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}
	var conf T
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("error unmarshal config file: %w", err)
	}
	return &conf, nil
}
