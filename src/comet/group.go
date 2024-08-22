package comet

type Group struct {
	id     int32
	groups map[UserID]*Channel
}

func (m *Manager) NewGroup() *Group {
	return &Group{
		groups: make(map[UserID]*Channel),
	}
}

func (g *Group) PutChannel(channel *Channel) {
	g.groups[channel.id] = channel
}

func (g *Group) DelChannel(channel *Channel) {
	delete(g.groups, channel.id)
}
