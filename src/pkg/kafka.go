package pkg

import (
	"laneIM/src/config"
	"log"

	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Client sarama.SyncProducer
}

type KafkaComsumer struct {
	Client *kafka.Consumer
}

func NewKafkaProducer(conf config.KafkaProducer) *KafkaProducer {
	kp := &KafkaProducer{}
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(conf.Addr, config)
	if err != nil {
		log.Panicln("faild to created kafka producer:", err)
		return kp
	}
	kp.Client = producer
	log.Println("create kafka producer:", conf.Addr)
	return kp
}

func NewKafkaComsumer(conf config.KafkaComsumer) *KafkaComsumer {
	kc := &KafkaComsumer{}
	config := &kafka.ConfigMap{
		"bootstrap.servers": conf.Addr,  // Kafka 服务器地址
		"group.id":          "im-group", // 消费者组 ID
		"auto.offset.reset": "earliest", // 从最早的消息开始消费
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Panicln("failed to create kafka comsumer:", err)
	}
	kc.Client = c
	log.Println("create kafka comsumer:", conf.Addr)
	return kc
}
