package server

import "github.com/spf13/viper"

type ServerConfig struct {
	Domain    string `mapstructure:"domain"`
	Port      string `mapstructure:"port"`
	Protected bool   `mapstructure:"protected"`
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

type GoogleOauthConfig struct {
	ClientID     string   `mapstructure:"clientID"`
	ClientSecret string   `mapstructure:"clientSecret"`
	RedirectURL  string   `mapstructure:"redirectUrl"`
	Scopes       []string `mapstructure:"scopes"`
}

type Config struct {
	Server            ServerConfig      `mapstructure:"server"`
	DB                DBConfig          `mapstructure:"db"`
	Auth              AuthConfig        `mapstructure:"auth"`
	GoogleOauthConfig GoogleOauthConfig `mapstructure:"googleOauth"`
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
