package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type MySQLConfig struct {
	Username string `mapstructure:"username"`
	Password string
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
}

func LoadConfig() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file, %s", err)
	}

	var configuration Config
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	configuration.MySQL.Password = os.Getenv("MYSQL_PASSWORD")

	return configuration
}
