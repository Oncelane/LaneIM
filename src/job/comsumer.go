package job

import (
	"log"

	"github.com/IBM/sarama"
)

type MyConsumer struct {
	// 可以添加其他字段，如日志记录器、统计信息等

}

func (consumer *MyConsumer) Setup(session sarama.ConsumerGroupSession) error {
	// 在这里可以执行一些初始化操作，比如日志记录
	return nil
}

func (consumer *MyConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	// 在这里可以执行一些清理操作
	return nil
}

func (consumer *MyConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received message: %s from %s/%d\n", message.Value, message.Topic, message.Partition)
		// 处理消息...

		// 标记消息为已消费
		session.MarkMessage(message, "")
	}
	return nil
}
