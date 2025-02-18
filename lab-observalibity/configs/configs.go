package configs

import (
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

type conf struct {
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	PortServiceA             string `mapstructure:"PORT_SERVICE_A"`
	PortB                    string `mapstructure:"PORT_B"`
	ServiceNameA             string `mapstructure:"SERVICE_NAME_A"`
	ServiceNameB             string `mapstructure:"SERVICE_NAME_B"`
	ServiceBUrl              string `mapstructure:"SERVICE_B_URL"`
	ZipkinUrl                string `mapstructure:"ZIPKIN_URL"`
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

type Config struct {
	OTELTracer         trace.Tracer
	RequestNameOTEL    string
	ServiceBUrl        string
	ExternalCallMethod string
	Content            string
}
