package pkg_test

import (
	"context"
	"fmt"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"log"
	"testing"

	"github.com/IBM/sarama"
)

var n int = 30
var topic string = "laneIMTest"

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
		fmt.Printf("Received message: %s from %s/%d\n", message.Value, message.Topic, message.Partition)
		// 处理消息...

		// 标记消息为已消费
		session.MarkMessage(message, "")
	}
	return nil
}
func TestPC(t *testing.T) {
	// go TestSaramComsumer(t)
	go TestSaramGroupComsumer(t)
	go TestSaramGroupComsumer(t)
	go TestProducer(t)
	select {}
}

func TestProducer(t *testing.T) {
	conf := config.KafkaProducer{
		Addr: []string{"127.0.0.1:9092"},
	}
	producer := pkg.NewKafkaProducer(conf)
	defer producer.Client.Close()
	for i := 0; i < n; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("%d", i)),
		}
		_, _, err := producer.Client.SendMessage(msg)
		if err != nil {
			t.Error(err)
		}
		if (i+1)%10000 == 0 {
			fmt.Println("send message", msg)
		}

	}

	log.Println("send all messages:")

}
func TestAsyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.NoResponse
	// config.Producer.Retry.Max = 10
	// config.Producer.Return.Successes = true
	produce, err := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		t.Error(err)
		return
	}
	defer produce.Close()

	for i := 0; i < n; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("%d", i)),
		}
		produce.Input() <- msg
		if (i+1)%10000 == 0 {
			fmt.Println("send message", msg)
		}
	}
	log.Println("send all messages:")
}

func TestSaramComsumer(t *testing.T) {
	// Kafka broker address
	brokers := []string{"localhost:9092"}

	// Create a new Sarama config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new Sarama consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start a new partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Start a loop to process messages
	log.Println("start recving...")
	count := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			count++
			if count == 10000 {
				count = 0
			}
			fmt.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Interrupt received, shutting down...")
			return
		}
	}

}
func TestSaramGroupComsumer(t *testing.T) {
	// Kafka broker address
	brokers := []string{"localhost:9092"}

	// Create a new Sarama config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new Sarama consumer
	consumer, err := sarama.NewConsumerGroup(brokers, "jobTest", config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer consumer.Close()
	handler := &MyConsumer{}

	// 注意：Consume 方法会阻塞，直到接收到上下文取消信号或发生错误
	ctx := context.Background() // 或者使用带取消的上下文
	if err := consumer.Consume(ctx, []string{topic}, handler); err != nil {
		log.Fatalf("Error from consumer group: %v", err)
	}
	log.Println("comsumer exit")
}

// func TestKafkaComsumer(t *testing.T) {
// 	conf := config.KafkaComsumer{
// 		Addr: "127.0.0.1:9092",
// 	}
// 	comsumer := pkg.NewKafkaComsumer(conf)
// 	defer comsumer.Client.Close()
// 	err := comsumer.Client.Subscribe(topic, nil)
// 	if err != nil {
// 		t.Errorf("faild to subscribe to topic %s: %s", topic, err)
// 	}
// 	count := 0
// 	for {
// 		msg, err := comsumer.Client.ReadMessage(-1)
// 		if err != nil {
// 			log.Printf("error wihle receiving message: %s\n", err)
// 			continue
// 		}
// 		count++
// 		log.Println("recv:", string(msg.Value))
// 		// if count == 1000 {
// 		// 	count = 0
// 		// 	log.Println("recv:", string(msg.Value))
// 		// }
// 	}
// 	log.Println("receive all msg")
// }
