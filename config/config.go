package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           int       `json:"port"`
	AiConfig       *AiConfig `json:"aiconfig"`
	Redis          *RedisCfg `json:"redis"`
	CrossEndpoint  string    `json:"crossChainEndpoint"`
	SwapEndpoint   string    `json:"swapEndpoint"`
	ConfigEndpoint string    `json:"configEndpoint"`
}

type AiConfig struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
	APIKey   string `json:"apikey"`
}

type RedisCfg struct {
	Addr         string        `json:"addr"`
	Password     string        `json:"password"`
	DB           int           `json:"db"`
	MinIdle      int           `json:"min_idle"`
	PoolSize     int           `json:"pool_size"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout"`
}

func LoadConfig(cfg interface{}) error {
	// Read in from .env file if available
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Load Config error: %s", err)
	}

	// Read in from environment variables
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("CROSSENDPOINT")
	_ = viper.BindEnv("SWAPENDPOINT")
	_ = viper.BindEnv("CONFIGENDPOINT")
	_ = viper.BindEnv("AICONFIG.ENDPOINT")
	_ = viper.BindEnv("AICONFIG.MODEL")
	_ = viper.BindEnv("AICONFIG.APIKEY")
	_ = viper.BindEnv("REDIS.ADDR")
	_ = viper.BindEnv("REDIS.PASSWORD")

	return viper.Unmarshal(cfg)
}
