package commands

import (
	"context"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/mq"
)

// Consume is the command for consuming messages from a topic.
func Consume(app *cli.MultipleProgram) {
	app.Register("consume", &cli.Command{
		Name:  "consume",
		Usage: "the consumer of mq",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "topic",
				Usage:   "the topic to consume",
				EnvVars: []string{"TOPIC"},
				Value:   "default",
			},
			&cli.StringFlag{
				Name:    "group",
				Usage:   "the group to consume",
				EnvVars: []string{"GROUP"},
				Value:   "default",
			},
			&cli.StringFlag{
				Name:     "consumer",
				Usage:    "the consumer to consume",
				EnvVars:  []string{"CONSUMER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "start",
				Usage:   "the start of the stream",
				EnvVars: []string{"START"},
				Value:   "$",
			},
			&cli.IntFlag{
				Name:    "batch-size",
				Usage:   "the batch size of the stream",
				EnvVars: []string{"BATCH_SIZE"},
				Value:   1,
			},
			&cli.StringFlag{
				Name:     "redis-host",
				Usage:    "the redis host",
				EnvVars:  []string{"REDIS_HOST"},
				Required: true,
			},
			&cli.IntFlag{
				Name:    "redis-port",
				Usage:   "the redis port",
				EnvVars: []string{"REDIS_PORT"},
				Value:   6379,
			},
			&cli.StringFlag{
				Name:    "redis-username",
				Usage:   "the redis username",
				EnvVars: []string{"REDIS_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "redis-password",
				Usage:   "the redis password",
				EnvVars: []string{"REDIS_PASSWORD"},
			},
			&cli.IntFlag{
				Name:    "redis-db",
				Usage:   "the redis db",
				EnvVars: []string{"REDIS_DB"},
				Value:   0,
			},
		},
		Action: func(ctx *cli.Context) error {
			ps := mq.New(&mq.Config{
				RedisHost:     ctx.String("redis-host"),
				RedisPort:     ctx.Int("redis-port"),
				RedisUsername: ctx.String("redis-username"),
				RedisPassword: ctx.String("redis-password"),
				RedisDB:       ctx.Int("redis-db"),
			})

			return ps.Consume(
				context.TODO(),
				ctx.String("topic"),
				ctx.String("group"),
				ctx.String("consumer"),
				ctx.String("start"),
				ctx.Int("batch-size"),
				func(msg *mq.Message) error {
					logger.Infof("consume message: %s", string(msg.Body))
					return nil
				},
			)
		},
	})
}
