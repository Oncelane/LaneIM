package comet

import (
	"laneIM/proto/msg"
	"log"
	"sync"
)

type Room struct {
	id     int64
	chsMap sync.Map
	Online int64
}

func (m *Bucket) NewRoom(roomid int64) *Room {
	m.mu.RLock()
	if r, exist := m.rooms[roomid]; exist {
		m.mu.RUnlock()
		log.Println("room:", roomid, "already exist,can not new")
		return r
	}
	m.mu.RUnlock()
	r := &Room{
		id: roomid,
	}
	m.mu.Lock()
	if r, exist := m.rooms[roomid]; exist {
		m.mu.RUnlock()
		log.Println("room:", roomid, "already exist,can not new")
		return r
	} else {
		m.rooms[roomid] = r
	}
	m.mu.Unlock()
	log.Println("new room:", roomid)
	return r
}

func (g *Room) PutChannel(channel *Channel) {
	g.chsMap.Store(channel.id, channel)
	// g.chs[channel.id] = channel
}

func (g *Room) DelChannel(channel *Channel) {
	g.chsMap.Delete(channel.id)
}

func (g *Room) Send(m *msg.Msg) {
	g.chsMap.Range(func(key, value any) bool {
		ch, ok := value.(*Channel)
		if ok {
			ch.sendCh <- m
		}
		return true
	})
}
