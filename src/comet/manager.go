package comet

import (
	"laneIM/src/pkg"
	"sync"
)

type Manager struct {
	pool  pkg.MsgPool
	rooms map[int64]*Room
	nRoom int
	mu    sync.Mutex
}

func (m *Manager) PutRoom(g *Room) {
	m.mu.Lock()
	m.nRoom++
	m.rooms[g.id] = g
	m.mu.Unlock()
}

func (m *Manager) DelRoom(id int64) {
	m.mu.Lock()
	delete(m.rooms, id)
	m.nRoom--
	m.mu.Unlock()
}

func (m *Manager) GetRoom(id int64) *Room {
	return m.rooms[id]
}

func (m *Manager) PutChannel(id int64, c *Channel) {
	m.GetRoom(id).PutChannel(c)
}

func (m *Manager) DelChannel(id int64, c *Channel) {
	m.GetRoom(id).DelChannel(c)
}
