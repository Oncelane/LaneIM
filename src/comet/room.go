package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg/laneLog"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Room struct {
	id         int64
	chsMap     sync.Map
	fullChsMap sync.Map
	Online     int64
}

func (m *Bucket) NewRoom(roomid int64) *Room {

	m.mu.Lock()
	defer m.mu.Unlock()
	if r, exist := m.rooms[roomid]; exist {
		laneLog.Logger.Infoln("[server] room:", roomid, "already exist,can not new")
		return r
	}
	newRoom := &Room{
		id: roomid,
	}
	laneLog.Logger.Infoln("[server] new room:", roomid)
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

func (g *Room) PutFullChannel(channel *Channel) {
	g.Online++
	g.fullChsMap.Store(channel.id, channel)
	// g.chs[channel.id] = channel
}
func (g *Room) DelFullChannel(channel *Channel) {
	g.Online--
	g.fullChsMap.Delete(channel.id)
}
func (g *Room) DelChannelBatch(in []*BatchStructSetOffline) {
	g.Online -= int64(len(in))
	for _, c := range in {
		g.chsMap.Delete(c.ch.id)
	}
}

// func (g *Room) Send(m *msg.Msg) {
// 	g.chsMap.Range(func(key, value any) bool {
// 		ch, ok := value.(*Channel)
// 		if ok {
// 			if ch.done {
// 				g.DelChannel(ch)
// 				return true
// 			}
// 			// laneLog.Logger.Infoln("message enter ch.sendch", ch.id)
// 			ch.sendCh <- m
// 		}
// 		return true
// 	})
// }

func (g *Room) SendBatch(m *msg.MsgBatch) {
	g.fullChsMap.Range(func(key, value any) bool {
		ch, ok := value.(*Channel)
		if ok {

			if ch.done {
				g.DelFullChannel(ch)
				return true
			}
			// laneLog.Logger.Infoln("message enter ch.sendch pass3 ", key, ch.id, m.String())
			ch.sendCh <- m
		}
		return true
	})
	onlyCountM := &msg.MsgBatch{
		Msgs: make([]*msg.Msg, 1),
	}

	onlyc := msg.COnlyCountMessage{
		Roomid: g.id,
		Count:  int32(len(m.Msgs)),
	}
	data, _ := proto.Marshal(&onlyc)
	onlyCountM.Msgs[0] = &msg.Msg{
		Path: "onlyCount",
		Data: data,
	}
	g.chsMap.Range(func(key, value any) bool {
		ch, ok := value.(*Channel)
		if ok {

			if ch.done {
				g.DelFullChannel(ch)
				return true
			}
			// laneLog.Logger.Infoln("message enter ch.sendch pass3 ", key, ch.id, m.String())
			ch.sendCh <- onlyCountM
		}
		return true
	})
}
