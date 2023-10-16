package commands

import (
	"context"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/mq"
)

// Send is the command for sending a message to a topic.
func Send(app *cli.MultipleProgram) {
	app.Register("send", &cli.Command{
		Name:  "send",
		Usage: "the send of mq",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "topic",
				Usage:   "the topic to send",
				EnvVars: []string{"TOPIC"},
				Value:   "default",
			},
			&cli.StringFlag{
				Name:     "message",
				Usage:    "the message to send",
				EnvVars:  []string{"MESSAGE"},
				Required: true,
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

			return ps.Send(context.TODO(), &mq.Msg{
				Topic: ctx.String("topic"),
				Body:  []byte(ctx.String("message")),
			})
		},
	})
}
