package comet

import "log"

type Room struct {
	id     int64
	chs    map[int64]*Channel
	Online int64
}

func (m *Bucket) NewRoom(roomid int64) *Room {
	if r, exist := m.rooms[roomid]; exist {
		log.Println("room:", roomid, "already exist,can not new")
		return r
	}
	r := &Room{
		id:  roomid,
		chs: make(map[int64]*Channel),
	}
	m.rooms[roomid] = r
	log.Println("new room:", roomid)
	return r
}

func (g *Room) PutChannel(channel *Channel) {
	g.chs[channel.id] = channel
}

func (g *Room) DelChannel(channel *Channel) {
	delete(g.chs, channel.id)
}
