package mq

// Message ...
type Message struct {
	// ID 消息ID
	ID string

	// Topic 主题
	Topic string

	// Body 消息体
	Body []byte

	// Group 消费者组
	Group string

	// 消费者组里的消费者
	// 同一组里的消费者共享消息，也就是说同一组里的消费者不会重复消费消息
	// 不同组里的消费者会重复消费消息
	Consumer string
}
