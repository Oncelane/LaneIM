package comet

type Room struct {
	id     int64
	rooms  map[int64]*Channel
	Online int64
}

func (m *Bucket) NewRoom() *Room {
	return &Room{
		rooms: make(map[int64]*Channel),
	}
}

func (g *Room) PutChannel(channel *Channel) {
	g.rooms[channel.id] = channel
}

func (g *Room) DelChannel(channel *Channel) {
	delete(g.rooms, channel.id)
}
