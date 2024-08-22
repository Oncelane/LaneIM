package comet

import (
	"laneIM/src/pkg"
	"sync"
)

type Manager struct {
	pool   pkg.MsgPool
	groups map[int32]*Group
	nGroup int
	mu     sync.Mutex
}

func (m *Manager) PutGroup(g *Group) {
	m.mu.Lock()
	m.nGroup++
	m.groups[g.id] = g
	m.mu.Unlock()
}

func (m *Manager) DelGroup(id int32) {
	m.mu.Lock()
	m.DelGroup(id)
	m.nGroup--
	m.mu.Unlock()
}

func (m *Manager) GetGroup(id int32) *Group {
	return m.groups[id]
}

func (m *Manager) PutChannel(id int32, c *Channel) {
	m.GetGroup(id).PutChannel(c)
}

func (m *Manager) DelChannel(id int32, c *Channel) {
	m.GetGroup(id).DelChannel(c)
}
