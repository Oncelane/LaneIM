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

	m.mu.Lock()
	defer m.mu.Unlock()
	if r, exist := m.rooms[roomid]; exist {
		log.Println("room:", roomid, "already exist,can not new")
		return r
	}
	newRoom := &Room{
		id: roomid,
	}
	log.Println("new room:", roomid)
	m.rooms[roomid] = newRoom
	return newRoom

}

func (g *Room) PutChannel(channel *Channel) {
	g.Online++
	g.chsMap.Store(channel.id, channel)
	// g.chs[channel.id] = channel
}

func (g *Room) DelChannel(channel *Channel) {
	g.Online--
	g.chsMap.Delete(channel.id)
}

func (g *Room) Send(m *msg.Msg) {
	g.chsMap.Range(func(key, value any) bool {
		ch, ok := value.(*Channel)
		if ok {
			if ch.done {
				g.DelChannel(ch)
				return true
			}
			// log.Println("room msg success send")
			ch.sendCh <- m
		}
		return true
	})
}
