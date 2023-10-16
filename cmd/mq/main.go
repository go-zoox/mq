package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/mq/cmd/mq/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:  "mq",
		Usage: "mq is the mq producer - send / consumer - consume",
	})

	commands.Consume(app)
	commands.Send(app)

	app.Run()
}
