// Package config служит для парсинга и последующей работы с файлами конфигурации.
package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

var ErrEnvVarNotFound = errors.New("you need to set up environment variable first.\nThat variable should contain 'conf.yaml' file")
var ErrFileNotFound = errors.New("config file is missing")

// GetConfig парсит файл конфигурации.
func GetConfig(path string) (*viper.Viper, error) {
	conf := viper.New()

	_, ok := os.LookupEnv(path)
	if ok == false {
		return nil, ErrEnvVarNotFound
	}

	conf.SetConfigFile(os.Getenv(path))
	if err := conf.ReadInConfig(); err != nil {
		return nil, ErrFileNotFound
	}
	// logging config.
	conf.SetDefault("OutputType", "console")

	// server config.
	conf.SetDefault("Addr", "localhost:8081")
	conf.SetDefault("ReadTimeout", time.Second*10)
	conf.SetDefault("WriteTimeout", time.Second*5)
	conf.SetDefault("IdleTimeout", time.Second*30)

	// Cache config.
	conf.SetDefault("RedisAddr", "localhost:6379")
	conf.SetDefault("RedisUser", "")
	conf.SetDefault("RedisPassword", "")
	conf.SetDefault("MaxRetries", 3)
	conf.SetDefault("PoolSize", 10)

	return conf, nil
}
