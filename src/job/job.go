package job

import (
	"context"
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/allegro/bigcache"
)

// singleton
type Job struct {
	etcd          *pkg.EtcdClient
	redis         *pkg.RedisClient
	db            *sql.SqlDB
	kafkaComsumer sarama.ConsumerGroup
	conf          config.Job
	mu            sync.RWMutex
	buckets       []*Bucket
	comets        map[string]*CometClient
	cache         *bigcache.BigCache
	daoo          *dao.Dao
}

func NewJob(conf config.Job) *Job {
	e := pkg.NewEtcd(conf.Etcd)

	// connect to redis
	addrs := e.GetAddr("redis")
	laneLog.Logger.Infof("job starting...\n get redis addrs: %v", addrs)

	j := &Job{
		etcd:          e,
		redis:         pkg.NewRedisClient(conf.Redis),
		kafkaComsumer: pkg.NewKafkaGroupComsumer(conf.KafkaComsumer),
		conf:          conf,
		comets:        make(map[string]*CometClient),
		daoo:          dao.NewDao(conf.Mysql.BatchWriter),
		cache:         localCache.Cache(time.Minute),
	}

	j.db = sql.NewDB(conf.Mysql)

	j.NewBucket()

	// wathc comet
	go j.WatchComet()
	go j.RunGroupComsumer()
	return j
}

func (j *Job) WatchComet() {
	for {
		addrs := j.etcd.GetAddr("grpc:comet")
		// laneLog.Logger.Infoln("addrs:", addrs)
		remoteAddrs := make(map[string]struct{})
		for _, addr := range addrs {
			remoteAddrs[addr] = struct{}{}
			// connet to comet
			if _, exist := j.comets[addr]; exist {
				continue
			}

			// discovery comet
			laneLog.Logger.Infoln("discovery comet:", addr)
			j.NewComet(addr)
		}
		for addr, c := range j.comets {
			if _, exist := remoteAddrs[addr]; !exist {
				j.mu.Lock()
				delete(j.comets, c.addr)
				c.conn.Close()
				for range j.conf.CometRoutineSize {
					c.done <- struct{}{}
				}
				laneLog.Logger.Infoln("remove comet:", addr)
				j.mu.Unlock()
			}

		}

		time.Sleep(time.Second)
	}
}

func (j *Job) Close() {
	laneLog.Logger.Infoln("job exit", j.conf.Addr)
}

func (j *Job) RunGroupComsumer() {
	handler := &MyConsumer{
		job: j,
	}
	if err := j.kafkaComsumer.Consume(context.Background(), j.conf.KafkaComsumer.Topics, handler); err != nil {
		log.Fatalf("Error from consumer group: %v", err)
	}
	laneLog.Logger.Infoln("group comsumer exit")
}
