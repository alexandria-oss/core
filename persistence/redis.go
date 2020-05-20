package persistence

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alexandria-oss/core/config"
	"github.com/go-redis/redis/v7"
)

// NewRedisPool Obtain a Redis connection pool
func NewRedisPool(cfg *config.Kernel) (*redis.Client, func(), error) {
	db, err := strconv.Atoi(cfg.InMemory.Database)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Network:            cfg.InMemory.Network,
		Addr:               cfg.InMemory.Host + fmt.Sprintf(":%d", cfg.InMemory.Port),
		Dialer:             nil,
		OnConnect:          nil,
		Password:           cfg.InMemory.Password,
		DB:                 db,
		MaxRetries:         10,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        30 * time.Second,
		ReadTimeout:        15 * time.Second,
		WriteTimeout:       15 * time.Second,
		PoolSize:           100,
		MinIdleConns:       32,
		MaxConnAge:         0,
		PoolTimeout:        24 * time.Second,
		IdleTimeout:        30 * time.Second,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})

	cleanup := func() {
		if client != nil {
			_ = client.Close()
		}
	}

	err = client.Ping().Err()
	if err != nil {
		return nil, cleanup, nil
	}

	return client, cleanup, nil
}
