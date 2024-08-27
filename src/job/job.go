package job

import (
	"laneIM/src/pkg"
	"log"
	"time"
)

// singleton
type Job struct {
	bucket *Bucket
	etcd   *pkg.EtcdClient
	redis  *pkg.RedisClient
}

func NewJob() *Job {
	e := pkg.NewEtcd()

	addrs := e.GetAddr("redis")
	log.Println("job starting...\n get redis addrs: %v", addrs)
	j := &Job{
		bucket: NewBucket(),
		etcd:   e,
		redis:  pkg.NewRedisClient(addrs),
	}

	// etcd
	go j.WatchComet()
	return j
}

func (j *Job) WatchComet() {
	for {
		addrs := j.etcd.GetAddr("comet")
		newComets := make(map[string]*ClientComet)
		for _, addr := range addrs {
			// connet to comet
			newComets[addr] = NewComet(addr)
		}
		j.bucket.rw.Lock()
		j.bucket.comets = newComets
		j.bucket.rw.Unlock()
		time.Sleep(time.Second)
	}
}
