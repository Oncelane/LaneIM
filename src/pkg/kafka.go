package pkg

import "log"

type KafkaProducer struct {
}

type KafkaComsumer struct {
}

func NewKafka(addr string) *KafkaProducer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 服务器地址
	}

	// 创建生产者
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
}
