package job

import (
	"context"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// singleton
type Job struct {
	etcd          *pkg.EtcdClient
	redis         *pkg.RedisClient
	kafkaComsumer sarama.ConsumerGroup
	conf          config.Job
	mu            sync.RWMutex
	buckets       []*Bucket
	comets        map[string]*CometClient
}

func NewJob(conf config.Job) *Job {
	e := pkg.NewEtcd(conf.Etcd)

	// connect to redis
	addrs := e.GetAddr("redis")
	log.Printf("job starting...\n get redis addrs: %v", addrs)
	j := &Job{
		etcd:          e,
		redis:         pkg.NewRedisClient(addrs),
		kafkaComsumer: pkg.NewKafkaGroupComsumer(conf.KafkaComsumer),
		conf:          conf,
		comets:        make(map[string]*CometClient),
	}

	j.NewBucket()

	// wathc comet
	go j.WatchComet()
	go j.RunGroupComsumer()
	return j
}

func (j *Job) WatchComet() {
	for {
		addrs := j.etcd.GetAddr("grpc:comet")
		for _, addr := range addrs {
			// connet to comet
			if _, exist := j.comets[addr]; exist {
				continue
			}
			log.Println("发现comet:", addr)
			j.NewComet(addr)
		}
		time.Sleep(time.Second)
	}
}

// func (j *Job) RunComsumer() {
// 	defer j.kafkaComsumer.Client.Close()
// 	// Start a new partition consumer
// 	topic := "laneIM"
// 	partitionConsumer, err := j.kafkaComsumer.Client.ConsumePartition(topic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		log.Fatalf("Failed to start partition consumer: %v", err)
// 	}
// 	defer partitionConsumer.Close()

// 	// Start a loop to process messages
// 	log.Println("start comsuming...")
// 	for {
// 		select {
// 		case msg := <-partitionConsumer.Messages():
// 			protoMsg := &logic.SendMsgReq{}
// 			err := proto.Unmarshal(msg.Value, protoMsg)
// 			if err != nil {
// 				log.Println("faild to ummarshal kafka msg:", err)
// 			}
// 			log.Printf("Received message: %s\n", protoMsg.String())

// 			// 从redis查询房间内其余comet
// 			cometAddr, err := model.RoomQueryComet(j.redis.Client, common.Int64(protoMsg.Roomid))
// 			for _, addr := range cometAddr {
// 				if cometClient, exist := j.bucket.comets[addr]; exist {
// 					cometClient.roomCh <- &comet.RoomReq{
// 						Roomid: protoMsg.Roomid,
// 						Data:   protoMsg.Data,
// 					}
// 				} else {
// 					log.Println("faild to find addr:", addr)
// 				}
// 			}

// 			// 获取comet地址

// 		case err := <-partitionConsumer.Errors():
// 			log.Printf("Error: %v\n", err)
// 			log.Println("Interrupt received, shutting down...")
// 			return
// 		}
// 	}
// }

func (j *Job) Close() {
	log.Println("job close:", j.conf.Addr)
}

func (j *Job) RunGroupComsumer() {
	handler := &MyConsumer{
		job: j,
	}
	if err := j.kafkaComsumer.Consume(context.Background(), j.conf.KafkaComsumer.Topics, handler); err != nil {
		log.Fatalf("Error from consumer group: %v", err)
	}
	log.Println("group comsumer exit")
}
