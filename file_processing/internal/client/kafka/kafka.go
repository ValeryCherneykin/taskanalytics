package kafka

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
