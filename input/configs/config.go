package configs

import "github.com/spf13/viper"

type conf struct {
	OrchestratorClientUrl    string `mapstructure:"ORCHESTRATOR_CLIENT_URL"`
	WebServerPort            string `mapstructure:"WEB_SERVER_PORT"`
	OtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelExporterZipkinUrl    string `mapstructure:"OTEL_EXPORTER_ZIPKIN_URL"`
	RequestNameOTEL          string `mapstructure:"REQUEST_NAME_OTEL"`
	OrchestratorSpanNameOTEL string `mapstructure:"ORCHESTRATOR_SPAN_NAME_OTEL"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(path + "/.env")
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
