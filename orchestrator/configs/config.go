package configs

import "github.com/spf13/viper"

type conf struct {
	LocationClientUrl        string `mapstructure:"LOCATION_CLIENT_URL"`
	WeatherClientUrl         string `mapstructure:"WEATHER_CLIENT_URL"`
	WeatherClientKey         string `mapstructure:"WEATHER_CLIENT_KEY"`
	WebServerPort            string `mapstructure:"WEB_SERVER_PORT"`
	OtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelExporterZipkinUrl    string `mapstructure:"OTEL_EXPORTER_ZIPKIN_URL"`
	RequestNameOTEL          string `mapstructure:"REQUEST_NAME_OTEL"`
	LocationSpanNameOTEL     string `mapstructure:"LOCATION_SPAN_NAME_OTEL"`
	WeatherSpanNameOTEL      string `mapstructure:"WEATHER_SPAN_NAME_OTEL"`
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
