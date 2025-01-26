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
	ServerAddress string
	Database      struct {
		WriteDSN string
		ReadDSN  string
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

	// Grpc Server.
	cfg.ServerAddress = viper.GetString("SERVER_ADDRESS")

	// Database.
	cfg.Database.WriteDSN = viper.GetString("DB_WRITE_DSN")
	cfg.Database.ReadDSN = viper.GetString("DB_READ_DSN")
}
