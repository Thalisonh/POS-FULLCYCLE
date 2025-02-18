package configs

import "github.com/spf13/viper"

type conf struct {
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	HostPortServiceA         string `mapstructure:"HOST_PORT_SERVICE_A"`
	HostPortServiceB         string `mapstructure:"HOST_PORT_SERVICE_B"`
	ServiceNameA             string `mapstructure:"SERVICE_NAME_A"`
	ServiceNameB             string `mapstructure:"SERVICE_NAME_B"`
}

func LoadConfig(path string) (*conf, error) {
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
	return cfg, err
}
