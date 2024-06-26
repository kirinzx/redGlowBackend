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
    Server HTTPServer `yaml:"http_server"`
    FrontServer FrontendServer `yaml:"frontend_server"`
    Postgres PostgresDB `yaml:"postgres"`
    Redis RedisDB `yaml:"redis"`
    AuthSettings Auth `yaml:"auth"`
    EmailSettings Email `yaml:"email"`
}

type HTTPServer struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
    Timeout     time.Duration `yaml:"timeout"`
    IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type FrontendServer struct {
    Host string `yaml:"host"`
    Scheme string `yaml:"scheme"`
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

type Auth struct {
    SessionExpiration time.Duration `yaml:"session_expiration"`
    SessionCookieName string `yaml:"session_cookie_name"`
    CSRFTokenCookiename string `yaml:"csrftoken_cookie_name"`
    CSRFTokenHeaderName string `yaml:"csrftoken_header_name"`
    UserSessionContextKey string `yaml:"user_session_context_key"`
    CodeExpiration time.Duration `yaml:"code_expiration"`
    FrontSignUp string `yaml:"sign_up_confirm_front_uri"`
    FrontPassRecovery string `yaml:"recovery_password_front_uri"`
}

type Email struct {
    HostUser string `yaml:"host_user"`
    Host string `yaml:"host"`
    HostPassword string `yaml:"host_password"`
    Port int `yaml:"port"`
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