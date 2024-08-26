package job

import (
	"laneIM/proto/comet"
	"sync"
)

type Bucket struct {
	rw     sync.RWMutex
	comets map[string]*Comet
	rooms  map[int64]*Room //roomid to room

}

func NewBucket() *Bucket {
	return &Bucket{
		comets: make(map[string]*Comet),
	}
}

// receive control msg

// message send to comet

func (g *Bucket) Brodcast(m *comet.BrodcastReq) {

	for _, c := range g.comets {
		c.brodcastCh <- m
	}
}

func (g *Bucket) Room(m *comet.RoomReq) {
	g.rooms[m.Roomid].Push(m)
}

func (g *Bucket) Single(m *comet.SingleReq) {
	g.rooms[m.Roomid].PushSingle(m)
}
