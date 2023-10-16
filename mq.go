package mq

// referrer:
//	https://juejin.cn/post/7058699128284381221
//	https://github.com/jiaxwu/rmq/blob/main/stream.go

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var errBusyGroup = "BUSYGROUP Consumer Group name already exists"

// MQ is a simple message queue.
type MQ interface {
	Send(ctx context.Context, msg *Message) error
	Consume(ctx context.Context, topic, group, consumer, start string, batchSize int, h Handler) error
}

// Handler is a mq handler.
type Handler func(msg *Message) error

type mq struct {
	client *redis.Client

	//
	maxLen int64
	approx bool
}

// Config is the config for a mq.
type Config struct {
	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	RedisDB       int

	//
	Redis *redis.Client

	// MaxLen 最大消息数量，如果大于这个数量，旧消息会被删除，0表示不管
	MaxLen int64
	// Approx 配合 MaxLen 使用的，表示几乎精确的删除消息，也就是不完全精确，由于stream内部是流，所以设置此参数xadd会更加高效
	Approx bool
}

// New creates a new mq.
func New(cfg *Config) MQ {
	client := cfg.Redis
	if client == nil {
		if cfg.RedisHost == "" {
			panic("redis config is empty")
		}

		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Username: cfg.RedisUsername,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
	}

	return &mq{
		client: client,
		//
		maxLen: cfg.MaxLen,
		approx: cfg.Approx,
	}
}

// Send 发送消息
func (q *mq) Send(ctx context.Context, msg *Message) error {
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: msg.Topic,
		MaxLen: q.maxLen,
		Approx: q.approx,
		ID:     "*",
		Values: []interface{}{"body", msg.Body},
	}).Err()
}

// Consume 返回值代表消费过程中遇到的无法处理的错误
// group 消费者组
// consumer 消费者组里的消费者
// batchSize 每次批量获取一批的大小
// start 用于创建消费者组的时候指定起始消费ID，0表示从头开始消费，$表示从最后一条消息开始消费
func (q *mq) Consume(ctx context.Context, topic, group, consumer, start string, batchSize int, h Handler) error {
	err := q.client.XGroupCreateMkStream(ctx, topic, group, start).Err()
	if err != nil && err.Error() != errBusyGroup {
		return err
	}
	for {
		// id = > 表示拉取新消息
		if err := q.consume(ctx, topic, group, consumer, ">", batchSize, h); err != nil {
			return err
		}

		// id = 0 表示拉取已经投递却未被ACK的消息，保证消息至少被成功消费1次
		if err := q.consume(ctx, topic, group, consumer, "0", batchSize, h); err != nil {
			return err
		}
	}
}

func (q *mq) consume(ctx context.Context, topic, group, consumer, id string, batchSize int, h Handler) error {
	// 阻塞的获取消息
	result, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		// Group 消费者组
		Group: group,
		// 消费者组里的消费者
		// 同一组里的消费者共享消息，也就是说同一组里的消费者不会重复消费消息
		// 不同组里的消费者会重复消费消息
		Consumer: consumer,
		// Streams 表示要消费的消息流
		// 这里后面还有一个">"其实是属于ID参数，表示只接收未投递给其他消费者的消息
		// 如果指定ID为数值，则表示只接收大于这个ID的已经被拉取却没有被ACK的消息
		// 所以我们这里先使用>拉取一次最新消息，再使用0拉取已经投递却没有ACK的消息，保证消息都能够成功消费
		Streams: []string{topic, id},
		// Count 表示一次性获取多少条消息，减少网络开销
		Count: int64(batchSize),
	}).Result()
	if err != nil {
		return err
	}

	// 处理消息
	for _, msg := range result[0].Messages {
		err := h(&Message{
			ID:       msg.ID,
			Topic:    topic,
			Body:     []byte(msg.Values["body"].(string)),
			Group:    group,
			Consumer: consumer,
		})
		if err == nil {
			err := q.client.XAck(ctx, topic, group, msg.ID).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
