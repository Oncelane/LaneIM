package comet

type Room struct {
	id     int64
	rooms  map[UserID]*Channel
	Online int64
}

func (m *Manager) NewRoom() *Room {
	return &Room{
		rooms: make(map[UserID]*Channel),
	}
}

func (g *Room) PutChannel(channel *Channel) {
	g.rooms[channel.id] = channel
}

func (g *Room) DelChannel(channel *Channel) {
	delete(g.rooms, channel.id)
}
