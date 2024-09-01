package comet

import (
	"sync"
)

type Bucket struct {
	rooms map[int64]*Room
	nRoom int
	mu    sync.Mutex
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

func (m *Bucket) GetRoom(id int64) *Room {
	return m.rooms[id]
}

func (m *Bucket) PutChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.PutChannel(c)
	}
}

func (m *Bucket) DelChannel(roomid int64, c *Channel) {
	if room := m.GetRoom(roomid); room != nil {
		room.DelChannel(c)
	}
}
