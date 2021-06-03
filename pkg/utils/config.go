package utils

import "github.com/spf13/viper"

type Config struct {
	Db   DatabaseConfig `mapstructure:",squash"`
	Auth AuthConfig     `mapstructure:",squash"`
}

type DatabaseConfig struct {
	InMemory    bool   `mapstructure:"DB_IN_MEMORY"`
	User        string `mapstructure:"DB_USER"`
	Password    string `mapstructure:"DB_PASSWORD"`
	Hostname    string `mapstructure:"DB_HOST"`
	Port        string `mapstructure:"DB_PORT"`
	Database    string `mapstructure:"DB_NAME"`
	DefaultUser string `mapstructure:"DB_DEFAULT_USER"`
	DefaultPw   string `mapstructure:"DB_DEFAULT_PASSWORD"`
}

type AuthConfig struct {
	Secret          string `mapstructure:"AUTH_SECRET"`
	Issuer          string `mapstructure:"AUTH_ISSUER"`
	ExpirationHours int64  `mapstructure:"AUTH_EXP"`
}

func LoadConfig(path string) (*Config, error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)

	return &config, err
}
