package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"

	"github.com/gorilla/websocket"
)

type Channel struct {
	id     UserID
	ws     pkg.RWMsg
	recvCh chan *msg.Msg
	sendCh chan *msg.Msg
}

func (m *Manager) newChannel(userId UserID, conn *websocket.Conn) *Channel {
	ch := &Channel{
		id:     userId,
		ws:     pkg.NewWs(conn, &m.pool),
		recvCh: make(chan *msg.Msg, 100),
		sendCh: make(chan *msg.Msg, 100),
	}
	ch.ServeIO()
	return ch
}

func (c *Channel) ServeIO() {
	go c.recvRoutine()
	go c.sendRoutine()
}

func (c *Channel) recvRoutine() {
	for {
		message, err := c.ws.ReadMsg()
		if err != nil {
			// TODO
			continue
		}
		c.recvCh <- message
	}
}

func (c *Channel) sendRoutine() {
	for message := range c.sendCh {
		err := c.ws.WriteMsg(message)
		if err != nil {
			// TODO
			continue
		}
	}
}
