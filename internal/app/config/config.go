package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
	DatabaseName         string `mapstructure:"DATABASE_NAME"`
	DatastoreMaxOpenConn *int32 `mapstructure:"DATABASE_MAX_OPEN_CONN"`
	DatastoreMinOpenConn *int32 `mapstructure:"DATABASE_MIN_OPEN_CONN"`
}

func (c *Config) IsDevelopmentMode() bool {
	return c.Environment == "development"
}

func (c *Config) IsTestMode() bool {
	return os.Getenv("TEST_MODE") == "1"
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
