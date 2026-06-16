package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"leap-one/service-task/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	CacheKeyPrefix = "leapone:task:"
	DefaultTTL     = 15 * time.Minute
)

// RedisClient Redis缓存客户端封装
type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

// InitRedis 初始化Redis连接并返回封装的客户端实�?
func InitRedis(cfg *config.RedisConfig, logger *zap.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) buildKey(key string) string {
	return CacheKeyPrefix + key
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, r.buildKey(key)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) GetString(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.buildKey(key)).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.SetWithTTL(ctx, key, value, DefaultTTL)
}

func (r *RedisClient) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存值失败: %w", err)
	}
	return r.client.Set(ctx, r.buildKey(key), data, ttl).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = r.buildKey(k)
	}
	return r.client.Del(ctx, redisKeys...).Err()
}

func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	redisKeys := make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = r.buildKey(k)
	}
	return r.client.Exists(ctx, redisKeys...).Result()
}
