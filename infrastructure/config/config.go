package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	SrvID         int      `yaml:"srvID"`
	MySQLAddr     string   `yaml:"mysqlAddr"`
	MySQLUser     string   `yaml:"mysqlUser"`
	MySQLPassword string   `yaml:"mysqlPassword"`
	MySQLDBName   string   `yaml:"mysqlDBName"`
	MySQLConLimit int      `yaml:"mysqlConLimit"`
	EtcdAddrs     []string `yaml:"etcdAddrs"`
	EtcdPath      string   `yaml:"etcdPath"`
	RedisAddr     string   `yaml:"redisAddr"`
	RedisPassword string   `yaml:redisPassword`
}

func NewConfig() (conf *Config, err error) {
	conf = &Config{}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	err = viper.Unmarshal(conf)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
