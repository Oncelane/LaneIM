package job

import (
	"laneIM/proto/comet"
	"log"
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
}

func (j *Job) Bucket(roomid int64) *Bucket {
	log.Println("choose bucket:", int(roomid)%len(j.buckets))
	return j.buckets[int(roomid)%len(j.buckets)]
}

func (b *Bucket) GetRoom(roomid int64) *Room {
	b.rw.RLock()
	if room, exist := b.rooms[roomid]; exist {
		b.rw.RUnlock()
		return room
	}
	b.rw.RUnlock()
	return b.NewRoom(roomid)
}

func (b *Bucket) NewRoom(roomid int64) *Room {

	// update from redis

	b.rw.Lock()
	defer b.rw.Unlock()
	if r, exist := b.rooms[roomid]; exist {
		log.Println("room:", roomid, "already exist,can not new")
		return r
	}
	newRoom := &Room{
		roomid: roomid,
	}
	log.Println("create new room:", roomid)
	newRoom.UpdateFromRedis(b.job.redis)
	log.Println("sync room from redis:", roomid)
	b.rooms[roomid] = newRoom

	return newRoom
}

// receive control msg

// message send to comet

func (j *Job) Brodcast(m *comet.BrodcastReq) {

	for _, c := range j.comets {
		c.brodcastCh <- m
	}
}

func (g *Bucket) Room(m *comet.RoomReq) {
	room := g.GetRoom(m.Roomid)
	room.rw.RLock()
	defer room.rw.RUnlock()
	for cometAddr := range room.info.Server {
		g.job.comets[cometAddr].roomCh <- m
	}

}

// func (g *Bucket) Single(m *comet.SingleReq) {
// 	g.rw.RLock()
// 	g.rooms[m.Roomid].PushSingle(m)
// 	g.rw.RUnlock()
// }
