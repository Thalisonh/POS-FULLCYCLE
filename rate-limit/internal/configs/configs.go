package configs

import (
	"strconv"

	"github.com/spf13/viper"
)

type conf struct {
	TokenLimit     string `mapstructure:"TOKEN_LIMIT"`
	IpLimit        string `mapstructure:"IP_LIMIT"`
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      string `mapstructure:"REDIS_PORT"`
	RateLimitIp    string `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken string `mapstructure:"RATE_LIMIT_TOKEN"`
	BlockTime      string `mapstructure:"BLOCK_TIME"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	rateLimitIP, _ := strconv.Atoi(cfg.RateLimitIp)
	rateLimitToken, _ := strconv.Atoi(cfg.RateLimitToken)
	blockTime, _ := strconv.Atoi(cfg.BlockTime)

	return &Config{
		RedisHost:      cfg.RedisHost,
		RedisPort:      cfg.RedisPort,
		RateLimitIP:    rateLimitIP,
		RateLimitToken: rateLimitToken,
		BlockTime:      blockTime,
	}, err
}

type Config struct {
	RedisHost      string
	RedisPort      string
	RateLimitIP    int
	RateLimitToken int
	BlockTime      int
}
