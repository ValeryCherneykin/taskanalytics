package kafka

import "context"

type Producer interface {
	SendMessage(ctx context.Context, topic string, key, value []byte) error
	Close() error
}

