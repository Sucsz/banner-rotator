package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort string `mapstructure:"http_port"`

	Postgres struct {
		Host     string        `mapstructure:"host"`
		Port     int           `mapstructure:"port"`
		User     string        `mapstructure:"user"`
		Password string        `mapstructure:"password"`
		DBName   string        `mapstructure:"dbname"`
		SSLMode  string        `mapstructure:"sslmode"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"postgres"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
	} `mapstructure:"kafka"`

	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig() (*Config, error) {
	// файл:
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// ENV: префикс и автоматическое чтение
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	// превращаем точки в ключах в подчёркивания для ENV:
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// дефолты:
	viper.SetDefault("http_port", "8080")
	viper.SetDefault("log_level", "info")

	viper.SetDefault("postgres.host", "postgres")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.dbname", "bannerdb")
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("postgres.timeout", 5*time.Second)

	viper.SetDefault("kafka.brokers", []string{"kafka:9092"})
	viper.SetDefault("kafka.topic", "banner-events")

	// читаем файл, но не падаем, если его нет, оставляем ENV или если ENV нет, то дефолты:
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("сonfig file not found: %v; falling back to ENV/defaults\n", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config into struct: %w", err)
	}
	return &cfg, nil
}
