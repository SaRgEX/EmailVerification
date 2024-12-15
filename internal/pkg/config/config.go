package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"time"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type HTTPServer struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Logger struct {
	FilePath string `yaml:"logger_path"`
	LogLevel string `yaml:"log_level"`
}

type SMTPServer struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	From     string `yaml:"from"`
	Password string
}

type Config struct {
	Database
	HTTPServer
	Logger
	SMTPServer
}

func New() *Config {
	return &Config{
		Database{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: envValue("DATABASE_PASSWORD", "postgres"),
			Database: viper.GetString("db.database"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
		HTTPServer{
			Addr:         viper.GetString("http_server.addr"),
			ReadTimeout:  viper.GetDuration("http_server.read_timeout"),
			WriteTimeout: viper.GetDuration("http_server.write_timeout"),
			IdleTimeout:  viper.GetDuration("http_server.idle_timeout"),
		},
		Logger{
			FilePath: viper.GetString("logger.file_path"),
			LogLevel: viper.GetString("logger.log_level"),
		},
		SMTPServer{
			Host:     viper.GetString("smtp_server.host"),
			Port:     viper.GetString("smtp_server.port"),
			From:     viper.GetString("smtp_server.from"),
			Password: envValue("SMTP_SERVER_PASSWORD", ""),
		},
	}
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.With("error", err).Error("Failed to load .env file")
		return nil
	}

	if err := ConfigureViper(); err != nil {
		slog.With("error", err).Error("Cannot read file")
		return nil
	}

	cfg := New()

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("error unmarshal configs: %s", err.Error())
	}

	return cfg
}

func ConfigureViper() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		slog.With("error", err).Error("Cannot read a config file")
		return err
	}

	return nil
}

func envValue(value, defaultValue string) string {
	if val := os.Getenv(value); val != "" {
		return val
	}
	return defaultValue
}
