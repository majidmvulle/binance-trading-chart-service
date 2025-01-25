package config_test

import (
	"github.com/majidmvulle/binance-trading-chart-service/ingestor/config"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	want := "golang-app-test"

	viper.SetDefault("APP_NAME", want)

	cfg := config.Config()

	if got := cfg.App.Name; got != want {
		t.Fatalf("config: expected %s, got: %s", want, got)
	}
}
