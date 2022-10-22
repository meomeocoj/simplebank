package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config store all configuration of application
// The values are read by viper from a config file or enviroment variable
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_Source"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecret         string        `mapstructure:"TOKEN_SECRET"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// loadConfig reads configuration from a yml or enviroment variable
func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
