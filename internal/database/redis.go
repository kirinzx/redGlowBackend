package database

import (
	"context"
	"encoding/json"
	"redGlow/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct{
	Client *redis.Client
}

func NewRedisDB(cfg *config.Config) *RedisDB{
	return &RedisDB{
		Client: redis.NewClient(&redis.Options{
			Addr: cfg.Redis.RedisAddr,
			Password: cfg.Redis.RedisPassword,
			Username: cfg.Redis.RedisUsername,
			DB: cfg.Redis.RedisDB,
			PoolFIFO: true,
			PoolSize: 25,
		}),
	}
}

func (rdb *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error{
	conn := rdb.Client.Conn()
	defer conn.Close()
	return conn.Set(ctx,key,value,expiration).Err()
}

func (rdb *RedisDB) Get(ctx context.Context, key string, typeToScan interface {}) error{
	conn := rdb.Client.Conn()
	defer conn.Close()
	value := conn.Get(ctx,key)
	if value.Err() != nil {
		return value.Err()
	}
	jsonBytes, err := value.Bytes()
	if err != nil{
		return err
	}
	err = json.Unmarshal(jsonBytes,typeToScan)

	if err != nil {
		return err
	}

	return nil
}

func (rdb *RedisDB) Del(ctx context.Context, key string) error{
	conn := rdb.Client.Conn()
	defer conn.Close()
	return conn.Del(ctx,key).Err()
}