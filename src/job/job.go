package job

import (
	"laneIM/src/pkg"
	"time"
)

// singleton
type Job struct {
	bucket *Bucket
	etcd   *pkg.EtcdClient
}

func NewJob() *Job {
	e := pkg.NewEtcd()

	j := &Job{
		bucket: NewBucket(),
		etcd:   e,
	}

	// etcd
	go j.WatchComet()
	return j
}

func (j *Job) WatchComet() {
	for {
		addrs := j.etcd.GetAddr("comet")
		newComets := make(map[string]*Comet)
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
