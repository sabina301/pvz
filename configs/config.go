package configs

import (
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

const (
	DefaultAPIRequestTimeout = time.Millisecond * 150
)

var AppConfiguration *AppConfig

type AppConfig struct {
	HttpPort int    `env:"HTTP_PORT" envDefault:"8080"`
	LogLvl   string `env:"LOG_LEVEL" envDefault:"info"`
	DB       DBConfig
	Auth     AuthConfig
}

type DBConfig struct {
	Host           string `env:"DB_HOST,required"`
	Port           int    `env:"DB_PORT" envDefault:"5432"`
	Name           string `env:"DB_NAME,required"`
	DriverName     string `env:"DB_DRIVER" envDefault:"postgres"`
	User           string `env:"DB_USER,required"`
	Password       string `env:"DB_PASSWORD,required"`
	MaxConnections int    `env:"MAX_CONNECTIONS" envDefault:"10"`
	SslMode        string `env:"DB_SSLMODE" envDefault:"disable"`
}

type AuthConfig struct {
	JwtSecret  string        `env:"JWT_SECRET,required"`
	Expiration time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()
	cfg := AppConfig{}
	cfg.DB = DBConfig{}
	cfg.Auth = AuthConfig{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	AppConfiguration = &cfg
	return &cfg, nil
}
