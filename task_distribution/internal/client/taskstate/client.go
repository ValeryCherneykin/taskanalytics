package taskstate

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	LPush(ctx context.Context, queue string, value interface{}) error
	BRPop(ctx context.Context, queue string, timeout time.Duration) (interface{}, error)
	Ping(ctx context.Context) error
	Close() error
}
