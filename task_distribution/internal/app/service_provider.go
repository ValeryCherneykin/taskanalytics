package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/api/taskqueue"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/kafka"
	producerPkg "github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/kafka/producer"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/taskstate"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/taskstate/redis"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository"
	taskstateRepo "github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository/taskqueue"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/service"
	taskstateService "github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/service/taskqueue"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	redisConfig config.RedisConfig
	grpcConfig  config.GRPCConfig
	kafkaConfig config.KafkaProducerConfig

	redisPool   *redigo.Pool
	redisClient taskstate.RedisClient

	taskQueueRepo repository.TaskQueue

	queueService service.QueueService

	taskQueue *taskqueue.Implementation

	kafkaProducer kafka.Producer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) KafkaConfig() config.KafkaProducerConfig {
	if s.kafkaConfig == nil {
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka config: %s", err.Error())
		}
		s.kafkaConfig = cfg
	}
	return s.kafkaConfig
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() taskstate.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) TaskQueueRepo(ctx context.Context) repository.TaskQueue {
	if s.taskQueueRepo == nil {
		s.taskQueueRepo = taskstateRepo.NewRepository(s.RedisClient())
	}
	return s.taskQueueRepo
}

func (s *serviceProvider) KafkaProducer() kafka.Producer {
	if s.kafkaProducer == nil {
		producer, err := sarama.NewSyncProducer(
			s.KafkaConfig().Brokers(),
			s.KafkaConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create kafka producer: %v", err)
		}

		s.kafkaProducer = producerPkg.NewKafkaProducer(producer)
	}
	return s.kafkaProducer
}

func (s *serviceProvider) QueueService(ctx context.Context) service.QueueService {
	if s.queueService == nil {
		s.queueService = taskstateService.NewService(s.TaskQueueRepo(ctx), s.KafkaProducer())
	}
	return s.queueService
}

func (s *serviceProvider) QueueImpl(ctx context.Context) *taskqueue.Implementation {
	if s.taskQueue == nil {
		s.taskQueue = taskqueue.NewImplementation(s.QueueService(ctx))
	}

	return s.taskQueue
}
