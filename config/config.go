package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		Jwt      JWT
		HTTP     HTTP
		Log      Log
		Pg       PG
		Security Security
		Email    Email
	}

	JWT struct {
		SecretKey       string        `env:"JWT_SECRET_KEY" env-required:"true"`
		AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" env-required:"true"`
		RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" env-required:"true"`
	}

	Email struct {
		FromMail     string        `env:"FROM_MAIL" env-required:"true"`
		MailPassword string        `env:"MAIL_PASSWORD" env-required:"true"`
		SmtpHost     string        `env:"SMTP_HOST" env-required:"true"`
		SmtpPort     int           `env:"SMTP_PORT" env-default:"465"`
		Timeout      time.Duration `env:"MAIL_TIMEOUT" env-default:"4m"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-required:"true"`
		Mode string `env:"GIN_MODE" env-default:"release"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" env-required:"true"`
		URL     string `env:"PG_URL" env-required:"true"`
	}

	Security struct {
		PasswordCost int `env:"SECURITY_PASSWORD_COST" env-default:"10"`
	}
)

func MustLoad() *Config {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("error loading .env file: %v", err)
		}
	}

	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	if cfg.Security.PasswordCost > 31 || cfg.Security.PasswordCost < 4 {
		log.Fatal("SECURITY_PASSWORD_COST is not allowed. It should be <31 or >=4")
	}

	return &cfg
}
