package redis

import (
	"context"
	"log"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/taskstate"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/config"
	"github.com/gomodule/redigo/redis"
)

var _ taskstate.RedisClient = (*client)(nil)

type handler func(ctx context.Context, conn redis.Conn) error

type client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

func NewClient(pool *redis.Pool, config config.RedisConfig) *client {
	return &client{
		pool:   pool,
		config: config,
	}
}

func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	return c.withConn(ctx, func(conn redis.Conn) error {
		_, err := conn.Do("SET", key, value)
		return err
	})
}

func (c *client) Get(ctx context.Context, key string) (interface{}, error) {
	var result interface{}
	err := c.withConn(ctx, func(conn redis.Conn) error {
		var err error
		result, err = conn.Do("GET", key)
		return err
	})
	return result, err
}

func (c *client) LPush(ctx context.Context, queue string, value interface{}) error {
	return c.withConn(ctx, func(conn redis.Conn) error {
		_, err := conn.Do("LPUSH", queue, value)
		return err
	})
}

func (c *client) BRPop(ctx context.Context, queue string, timeout time.Duration) (interface{}, error) {
	var result []interface{}
	err := c.withConn(ctx, func(conn redis.Conn) error {
		var err error
		result, err = redis.Values(conn.Do("BRPOP", queue, int(timeout.Seconds())))
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, nil
	}
	return result[1], nil
}

func (c *client) Ping(ctx context.Context) error {
	return c.withConn(ctx, func(conn redis.Conn) error {
		_, err := conn.Do("PING")
		return err
	})
}

func (c *client) Close() error {
	return c.pool.Close()
}

func (c *client) withConn(ctx context.Context, fn func(redis.Conn) error) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close redis connection: %v", cerr)
		}
	}()
	return fn(conn)
}
