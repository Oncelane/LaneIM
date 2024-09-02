package job

import (
	"laneIM/proto/comet"
	"sync"
)

type Bucket struct {
	rw     sync.RWMutex
	comets map[string]*CometClient
	rooms  map[int64]*Room //roomid to room
	job    *Job
}

func (j *Job) NewBucket() {
	j.buckets = make([]*Bucket, j.conf.BucketSize)
	for i := range j.conf.BucketSize {
		j.buckets[i] = &Bucket{
			comets: make(map[string]*CometClient),
			rooms:  make(map[int64]*Room),
			job:    j,
		}
	}
	return
}

func (j *Job) Bucket(roomid int64) *Bucket {
	return j.buckets[int(roomid)%len(j.buckets)]
}

// receive control msg

// message send to comet

func (j *Job) Brodcast(m *comet.BrodcastReq) {

	for _, c := range j.comets {
		c.brodcastCh <- m
	}
}

func (g *Bucket) Room(m *comet.RoomReq) {
	g.rw.RLock()
	if room, exist := g.rooms[m.Roomid]; exist {
		g.rw.RUnlock()
		room.rw.RLock()
		for cometAddr := range room.info.Server {
			g.job.comets[cometAddr].roomCh <- m
		}
		room.rw.RUnlock()
	} else {
		g.rw.RUnlock()
		room = g.job.NewRoom(m.Roomid)
		room.rw.RLock()
		for cometAddr := range room.info.Server {
			g.job.comets[cometAddr].roomCh <- m
		}
		room.rw.RUnlock()
	}

}

// func (g *Bucket) Single(m *comet.SingleReq) {
// 	g.rw.RLock()
// 	g.rooms[m.Roomid].PushSingle(m)
// 	g.rw.RUnlock()
// }
