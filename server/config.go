package server

import "github.com/spf13/viper"

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DbName   string `mapstructure:"dbname"`
	SslMode  string `mapstructure:"sslmode"`
	DbType   string `mapstructure:"dbtype"`
}

type AuthConfig struct {
	Key         string `mapstructure:"key"`
	JKey        string `mapstructure:"jKey"`
	ExpDuration int    `mapstructure:"expDuration"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("cfg")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}
