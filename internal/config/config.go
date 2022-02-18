package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/gophermart/pkg/logging"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	accrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	Key                  string `env:"KEY"`
	PasswordPepper       string `env:"PASSWORD_PEPPER"`
}

func NewConfig() Config {
	cfg := Config{
		Address:              "127.0.0.1:8080",
		accrualSystemAddress: "127.0.0.1:8081",
		Key:                  "SecretKey",
		PasswordPepper:       "pepper",
		DatabaseUri:          "postgres://localhost:5432/gophermart?sslmode=disable",
	}

	cfg.parseFlags()
	cfg.parseEnv()

	return cfg
}

func (cfg *Config) parseEnv() {
	err := env.Parse(cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msg("error parse environment")
	}
}

func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.StringVar(&cfg.DatabaseUri, "d", cfg.DatabaseUri, "Database uri")
	flag.StringVar(&cfg.accrualSystemAddress, "r", cfg.accrualSystemAddress, "accrual system address")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "jwt secret key")

	flag.Parse()
}
