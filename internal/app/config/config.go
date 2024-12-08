package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

const (
	mainEnvFile = ".env"
	hide        = "hide"
)

type Config struct {
	Environment string `mapstructure:"ENV"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	RestAPIPort string `mapstructure:"REST_API_PORT"`

	DatabaseURL          string `mapstructure:"DATABASE_URL"`
	DatastoreMaxOpenConn *int32 `mapstructure:"DATABASE_MAX_OPEN_CONN"`
	DatastoreMinOpenConn *int32 `mapstructure:"DATABASE_MIN_OPEN_CONN"`

	RebalancePoolThresholdPercent float64           `mapstructure:"REBALANCE_POOL_THRESHOLD_PERCENT"`
	RebalanceCheckInterval        time.Duration     `mapstructure:"REBALANCE_CHECK_INTERVAL"`
	CurrenciesEnabled             []domain.Currency `mapstructure:"CURRENCIES_ENABLED"`
}

func (c *Config) IsDevelopmentMode() bool {
	return c.Environment == "development"
}

func (c *Config) IsTestMode() bool {
	return os.Getenv("TEST_MODE") == "1"
}

func (c *Config) LogFields() map[string]any {
	return map[string]any{
		"Environment":                   c.Environment,
		"LogLevel":                      c.LogLevel,
		"RestAPIPort":                   c.RestAPIPort,
		"DatabaseURL":                   hide, // c.DatabaseURL,
		"DatastoreMaxOpenConn":          c.DatastoreMaxOpenConn,
		"DatastoreMinOpenConn":          c.DatastoreMinOpenConn,
		"RebalancePoolThresholdPercent": c.RebalancePoolThresholdPercent,
		"RebalanceCheckInterval":        c.RebalanceCheckInterval,
		"CurrenciesEnabled":             c.CurrenciesEnabled,
	}
}

func Load(filenames ...string) (*Config, error) {
	var cfg = &Config{}

	filenames = append(filenames, mainEnvFile)

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	for _, filename := range filenames {
		if _, err := os.Stat(filename); err != nil {
			continue
		}

		viper.SetConfigFile(filename)

		if err := viper.ReadInConfig(); err != nil {
			return nil, errors.Wrapf(err, "error to read config, path: %s", mainEnvFile)
		}

		if err := viper.MergeInConfig(); err != nil {
			return nil, errors.Wrapf(err, "error to merge config, filename: %s", filename)
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, errors.Wrapf(err, "error to unmarshal config, filename: %s", filename)
		}
	}

	return cfg, nil
}
