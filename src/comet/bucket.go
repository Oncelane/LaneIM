package comet

import (
	"sync"
)

type Bucket struct {
	rooms map[int64]*Room
	nRoom int
	mu    sync.RWMutex
}

func NewBucket() *Bucket {
	return &Bucket{
		rooms: make(map[int64]*Room),
		nRoom: 0,
	}
}

func (m *Bucket) PutRoom(g *Room) {
	m.mu.Lock()
	m.nRoom++
	m.rooms[g.id] = g
	m.mu.Unlock()
}

func (m *Bucket) DelRoom(roomid int64) {
	m.mu.Lock()
	delete(m.rooms, roomid)
	m.nRoom--
	m.mu.Unlock()
}

func (m *Bucket) GetRoom(roomid int64) *Room {
	m.mu.RLock()
	if r, exist := m.rooms[roomid]; exist {
		m.mu.RUnlock()
		return r
	}
	m.mu.RUnlock()
	r := m.NewRoom(roomid)
	return r
}

func (m *Bucket) PutChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.PutChannel(c)
	}
}
func (m *Bucket) PutFullChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.PutFullChannel(c)
	}
}

func (m *Bucket) DelFullChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.DelFullChannel(c)
	}
}

func (m *Bucket) DelChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.DelChannel(c)
		room.DelFullChannel(c)
	}
}

func (m *Bucket) DelChannelAll(c *Channel) {
	for _, room := range m.rooms {
		room.DelChannel(c)
		room.DelFullChannel(c)
	}
}
func (m *Bucket) DelChannelAllBatch(in []*BatchStructSetOffline) {
	for _, room := range m.rooms {
		room.DelChannelBatch(in)
	}
}
