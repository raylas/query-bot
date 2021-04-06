package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	ListQueries bool
	Queries     []Query
}

type Query struct {
	Command string `mapstructure:"command"`
	Method  string `mapstructure:"method"`
	URL     string `mapstructure:"url"`
	File    bool   `mapstructure:"file"`
}

func Load() (config Configuration, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath("./")
	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatalf("[ERROR] %s \n", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("[ERROR] %s \n", err)
	}

	return
}
