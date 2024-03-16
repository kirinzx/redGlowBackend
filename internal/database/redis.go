package database

import (
	"context"
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
			Addr: cfg.RedisDB.RedisAddr,
			Password: cfg.RedisDB.RedisPassword,
			Username: cfg.RedisDB.RedisUsername,
			DB: cfg.RedisDB.RedisDB,
		}),
	}
}

func (rdb *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error{
	conn := rdb.Client.Conn()
	defer conn.Close()
	return conn.Set(ctx,key,value,expiration).Err()
}

func (rdb *RedisDB) Get(ctx context.Context, key string, typeToScan interface{}) any{
	conn := rdb.Client.Conn()
	defer conn.Close()
	value := conn.Get(ctx,key)

	return value.Scan(typeToScan)
}

func (rdb *RedisDB) Del(ctx context.Context, key string) error{
	conn := rdb.Client.Conn()
	defer conn.Close()
	return conn.Del(ctx,key).Err()
}