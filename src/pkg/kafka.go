package pkg

import (
	"laneIM/src/config"
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Client sarama.SyncProducer
}

type KafkaComsumer struct {
	Client sarama.Consumer
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
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	c, err := sarama.NewConsumer(conf.Addr, config)
	if err != nil {
		log.Panicln("failed to create kafka comsumer:", err)
	}
	kc.Client = c
	log.Println("create kafka comsumer:", conf.Addr)
	return kc
}

func NewKafkaGroupComsumer(conf config.KafkaComsumer) sarama.ConsumerGroup {
	// Create a new Sarama config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new Sarama consumer
	consumer, err := sarama.NewConsumerGroup(conf.Addr, conf.GroupId, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	return consumer
}
