package redisclient

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type MyRedisClient struct {
	*redis.Client
}

func NewMyRedisClient() *MyRedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &MyRedisClient{client}
}

func (rc *MyRedisClient) SetKey(key, value string, expiration time.Duration) error {
	return rc.Set(context.Background(), key, value, expiration).Err()
}

func (rc *MyRedisClient) GetKey(key string) (string, error) {
	return rc.Get(context.Background(), key).Result()
}

func (rc *MyRedisClient) GetAll() ([]string, error) {
	return rc.Keys(context.Background(), "*").Result()
}

func (rc *MyRedisClient) GetHashValues(key string) ([]string, error) {
	return rc.HVals(context.Background(), key).Result()
}
