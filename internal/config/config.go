package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB struct {
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     string `env:"DB_PORT" env-default:"5432"`
		User     string `env:"DB_USER" env-default:"user"`
		Password string `env:"DB_PASSWORD" env-default:"password"`
		Name     string `env:"DB_NAME" env-default:"Price"`
		SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
	}

	Kafka struct {
		Brokers []string `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
		Topic   string   `env:"KAFKA_TOPIC" env-default:"price-updates"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" env-default:"localhost"`
		Port     string `env:"REDIS_PORT" env-default:"6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
	}

	Health struct {
		Port string `env:"PORT" env-default:":8080"`
	}

	App struct {
		GRPCServerPort string `env:"GRPC_SERVER_PORT" env-default:":50051"`
		ParserAddress  string `env:"PARSER_ADDRESS" env-default:"parser:50051"`
	}
}

func MustLoad() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("не удалось прочитать конфиг: %v", err)
	}
	return &cfg
}
