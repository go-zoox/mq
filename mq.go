package mq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// MQ is a simple message queue.
type MQ interface {
	Send(ctx context.Context, msg *Msg) error
	Consume(ctx context.Context, topic string, handler Handler) error
}

// Handler is a mq handler.
type Handler func(msg *Msg) error

type mq struct {
	client *redis.Client
}

// Config is the config for a mq.
type Config struct {
	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	RedisDB       int
}

// New creates a new mq.
func New(cfg *Config) MQ {
	return &mq{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Username: cfg.RedisUsername,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}),
	}
}

func (m *mq) Send(ctx context.Context, msg *Msg) error {
	return m.client.LPush(ctx, msg.Topic, msg.Body).Err()
}

func (m *mq) Consume(ctx context.Context, topic string, h Handler) error {
	for {
		// get message body
		body, err := m.client.LIndex(ctx, topic, -1).Bytes()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return err
			}

			time.Sleep(time.Second)
			continue
		}

		// handler message
		err = h(&Msg{
			Topic: topic,
			Body:  body,
		})
		if err != nil {
			continue
		}

		// remove message from queue, as ack
		err = m.client.RPop(ctx, topic).Err()
		if err != nil {
			return err
		}
	}
}
