package config

import (
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port	string
		Host	string
	}

	Database struct{
		Host	string
		Port	string
		User	string
		Password	string
		DBName	string
		SSLMode	string
	}

	JWT struct {
		Secret	string
		TokenExpiry time.Duration
	}

	RabbitMQ struct {
		Host     string
		Port     string
		User     string
		Password string
	}

	Enviroment string
}

func Load() (*Config, error){
	godotenv.Load()

	cfg := &Config{}

	//Server
	cfg.Server.Port = getEnv("SERVER_PORT","8082")	
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")

	//Database
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "admin1")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "auth_db")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	//JWT
	cfg.JWT.Secret = getEnv("JWT_SECRET", "")
	cfg.JWT.TokenExpiry = time.Hour * 24

	// RabbitMQ
	cfg.RabbitMQ.Host = getEnv("RABBITMQ_HOST", "localhost")
	cfg.RabbitMQ.Port = getEnv("RABBITMQ_PORT", "5672")
	cfg.RabbitMQ.User = getEnv("RABBITMQ_USER", "guest")
	cfg.RabbitMQ.Password = getEnv("RABBITMQ_PASSWORD", "guest")

	cfg.Enviroment = getEnv("ENV", "development")
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string{
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname =%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}
