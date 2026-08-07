package config

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	CORS     CORSConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port    int
	RunMode string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
}

type CORSConfig struct {
	AllowOrigins string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SSLMode  bool
}

type RedisConfig struct {
	Host               string
	Port               int
	Password           string
	Db                 int
	MinIdleConnections int
	PoolSize           int
	PoolTimeout        int
}

func getConfigpatch(env string) string {
	switch env {
	case "docker":
		return "config/config-docker"
	case "production":
		return "config/config-production"
	default:
		return "../config/config-development"
	}
}

func LoadConfig(filename, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}

func GetConfig() *Config {
	cfgPath := getConfigpatch(os.Getenv("APP_ENV"))
	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Error in Load config %v", err)
	}
	cfg, err := ParseConfig(v)
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}
