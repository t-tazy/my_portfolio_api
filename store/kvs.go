package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/t-tazy/my_portfolio_api/config"
	"github.com/t-tazy/my_portfolio_api/entity"
)

type KVS struct {
	Cli *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	// redis serverへのclientを返す
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

// アクセストークンのIDをキー、ユーザーIDを値として保存する
func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	id := int64(userID)
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
