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
	log.Printf("%p new room:%d before:%p", m, roomid, m.rooms[roomid])
	m.rooms[roomid] = newRoom
	log.Printf("%p new room:%d after:%p", m, roomid, m.rooms[roomid])
	return newRoom

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
			log.Println("room msg success send")
			ch.sendCh <- m
		}
		return true
	})
}
