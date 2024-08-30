package job

import (
	"laneIM/proto/comet"
	"laneIM/proto/logic"
	"laneIM/src/common"
	"laneIM/src/config"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

// singleton
type Job struct {
	bucket        *Bucket
	etcd          *pkg.EtcdClient
	redis         *pkg.RedisClient
	kafkaComsumer *pkg.KafkaComsumer
	conf          config.Job
}

func NewJob(conf config.Job) *Job {
	e := pkg.NewEtcd(conf.Etcd)

	// connect to redis
	addrs := e.GetAddr("redis")
	log.Printf("job starting...\n get redis addrs: %v", addrs)
	j := &Job{
		bucket: NewBucket(),
		etcd:   e,
		redis:  pkg.NewRedisClient(addrs),
		conf:   conf,
	}

	// connet to kafka
	j.kafkaComsumer = pkg.NewKafkaComsumer(conf.KafkaComsumer)

	// wathc comet
	go j.WatchComet()
	go j.RunComsumer()
	return j
}

func (j *Job) WatchComet() {
	for {
		addrs := j.etcd.GetAddr("grpc:comet")
		for _, addr := range addrs {
			// connet to comet
			if _, exist := j.bucket.comets[addr]; exist {
				continue
			}
			log.Println("发现comet:", addr)
			j.bucket.comets[addr] = NewComet(addr)
		}

		time.Sleep(time.Second)
	}
}

func (j *Job) RunComsumer() {
	defer j.kafkaComsumer.Client.Close()
	// Start a new partition consumer
	topic := "laneIM"
	partitionConsumer, err := j.kafkaComsumer.Client.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Start a loop to process messages
	log.Println("start comsuming...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			protoMsg := &logic.SendMsgReq{}
			err := proto.Unmarshal(msg.Value, protoMsg)
			if err != nil {
				log.Println("faild to ummarshal kafka msg:", err)
			}
			log.Printf("Received message: %s\n", protoMsg.String())

			// 从redis查询房间内其余comet
			cometAddr, err := model.RoomQueryComet(j.redis.Client, common.Int64(protoMsg.Roomid))
			for _, addr := range cometAddr {
				if cometClient, exist := j.bucket.comets[addr]; exist {
					cometClient.roomCh <- &comet.RoomReq{
						Roomid: protoMsg.Roomid,
						Data:   protoMsg.Data,
					}
				} else {
					log.Println("faild to find addr:", addr)
				}
			}

			// 获取comet地址

		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v\n", err)
			log.Println("Interrupt received, shutting down...")
			return
		}
	}
}

func (j *Job) Close() {
	log.Println("job close:", j.conf.Addr)
}
