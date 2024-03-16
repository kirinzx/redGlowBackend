package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
    HTTPServer `yaml:"http_server"`
    PostgresDB `yaml:"postgres"`
    RedisDB `yaml:"redis"`
}

type HTTPServer struct {
    Address     string        `yaml:"address" env-default:"0.0.0.0:8000"`
    Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresDB struct {
    PostgresHost string `yaml:"host"`
    PostgresPort string `yaml:"port"`
    PostgresUser string `yaml:"user"`
    PostgresPassword string `yaml:"password"`
    PostgresDatabaseName string `yaml:"db_name"`
}

type RedisDB struct {
    RedisAddr string `yaml:"addr"`
    RedisPassword string `yaml:"password"`
    RedisDB int `yaml:"db"`
    RedisUsername string `yaml:"username"`
}

func NewConfig(logger *zap.Logger) *Config {
    if err := godotenv.Load(); err != nil {
        logger.Fatal("No .env file found")
    }
    
    path := os.Getenv("CONFIG_PATH")
    if path == "" {
        logger.Fatal("CONFIG_PATH environment variable is not set")
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        logger.Fatal("Error reading file")
    }
    replaced := os.ExpandEnv(string(data))
    cfg := &Config{}
    err = yaml.Unmarshal([]byte(replaced), cfg)
    
    if err != nil{
        logger.Fatal(fmt.Sprintf("Error creating Config. Error: %s",err))
    }
    return cfg
}