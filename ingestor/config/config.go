package config

import (
	"github.com/spf13/viper"
)

var cfg *AppConfig

type AppConfig struct {
	App struct {
		Name  string
		Debug bool
		Env   string
	}

	Binance struct {
		WebsocketBaseURL string
		Symbols          []string
	}
}

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg = &AppConfig{}

	// App.
	cfg.App.Name = viper.GetString("APP_NAME")
	cfg.App.Debug = viper.GetBool("APP_DEBUG")
	cfg.App.Env = viper.GetString("APP_ENV")

	// Binance.
	cfg.Binance.WebsocketBaseURL = viper.GetString("BINANCE_WEBSOCKET_BASE_URL")
	cfg.Binance.Symbols = viper.GetStringSlice("BINANCE_SYMBOLS")
}
