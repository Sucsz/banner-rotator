// Package config содержит структуры для чтения .env и флагов.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// PostgresConfig описывает настройки подключения к PostgreSQL.
type PostgresConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	DBName   string        `mapstructure:"dbname"`
	SSLMode  string        `mapstructure:"sslmode"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// KafkaConfig описывает параметры подключения к Kafka.
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

// Config основная структура конфигурации приложения.
type Config struct {
	HTTPPort string         `mapstructure:"http_port"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	LogLevel string         `mapstructure:"log_level"`
	Epsilon  float64        `mapstructure:"epsilon"`
}

// LoadConfig загружает конфигурацию: сначала defaults и файл, затем ENV-override.
func LoadConfig() (*Config, error) {
	// 1) Значения по умолчанию
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

	viper.SetDefault("epsilon", 0.1)

	// 2) Чтение файла конфигурации (приоритет над defaults)
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("config file not found: %v; falling back to defaults/ENV\n", err)
	}

	// 3) ENV-override (самый высокий приоритет)
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 4) Маппинг значений в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config into struct: %w", err)
	}
	return &cfg, nil
}
