package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Env        string `env:"PORTFOLIO_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"PORTFOLIO_DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"PORTFOLIO_DB_PORT" envDefault:"33306"`
	DBUser     string `env:"PORTFOLIO_DB_USER" envDefault:"portfolio"`
	DBPassword string `env:"PORTFOLIO_DB_PASSWORD" envDefault:"portfolio"`
	DBName     string `env:"PORTFOLIO_DB_NAME" envDefault:"portfolio"`
	RedisHost  string `env:"PORTFOLIO_REDIS_HOST" envdefault:"127.0.0.1"`
	RedisPort  int    `env:"PORTFOLIO_REDIS_PORT" envdefault:"36379"`
}

// 環境変数から情報を取得
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
