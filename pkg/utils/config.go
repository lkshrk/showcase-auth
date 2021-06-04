package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
}

type DatabaseConfig struct {
	InMemory        bool `mapstructure:"in-memory"`
	User            string
	Password        string
	Hostname        string
	Port            string
	DB              string
	DefaultUser     string `mapstructure:"default-user"`
	DefaultPassword string `mapstructure:"default-password"`
}

type AuthConfig struct {
	Secret          string
	Issuer          string
	ExpirationHours int64 `mapstructure:"expiration-hours"`
}

func LoadConfig() (*Config, error) {

	configPath, isFound := os.LookupEnv("CONFIG_PATH")
	if !isFound {
		// attempt to load config at baseDir
		configPath = "../../"
	}
	configName, isFound := os.LookupEnv("CONFIG_NAME")
	if !isFound {
		configName = "config"
	}
	configType, isFound := os.LookupEnv("CONFIG_TYPE")
	if !isFound {
		configType = "yml"
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error reading %s%s.%s file, only reading config from ENV\n", configPath, configName, configType)
	}

	var config Config
	err = viper.Unmarshal(&config)

	return &config, err
}
