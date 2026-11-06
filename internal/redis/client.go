package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/config"
)

func Open(ctx context.Context, cfg config.Config) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
