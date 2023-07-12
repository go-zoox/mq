# MQ - lightweight message queue

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/mq)](https://pkg.go.dev/github.com/go-zoox/mq)
[![Build Status](https://github.com/go-zoox/mq/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/mq/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/mq)](https://goreportcard.com/report/github.com/go-zoox/mq)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/mq/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/mq?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/mq.svg)](https://github.com/go-zoox/mq/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/mq.svg?label=Release)](https://github.com/go-zoox/mq/tags)

## Installation
To install the package, run:
```bash
go get github.com/go-zoox/mq
```

## Getting Started


### Consumer
```go
import (
  "github.com/go-zoox/mq"
)

func main(t *testing.T) {
	m := mq.New(&mq.Config{
		RedisHost:     <RedisHost>,
		RedisPort:     <RedisPort>,
		RedisUsername: <RedisUsername>,
		RedisPassword: <RedisPassword>,
		RedisDB:       <RedisDB>,
	})

	m.Consume(context.TODO(), "default", func(msg *mq.Message) error {
		logger.Infof("received message: %s", string(msg.Body))
		return nil
	})
}
```

### Producer
```go
import (
  "github.com/go-zoox/mq"
)

func main(t *testing.T) {
	m := mq.New(&mq.Config{
		RedisHost:     <RedisHost>,
		RedisPort:     <RedisPort>,
		RedisUsername: <RedisUsername>,
		RedisPassword: <RedisPassword>,
		RedisDB:       <RedisDB>,
	})

	m.Send(context.TODO(), &mq.Message{
		Topic: "default",
		Body:  []byte("hello world"),
	})
}
```

## License
GoZoox is released under the [MIT License](./LICENSE).
