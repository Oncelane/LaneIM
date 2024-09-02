package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"

	"github.com/gorilla/websocket"
)

type Channel struct {
	id   int64
	conn pkg.MsgReadWriteCloser
	// conn   pkg.MsgReadWriteCloser
	recvCh chan *msg.Msg
	sendCh chan *msg.Msg
	done   bool
}

func (c *Comet) NewChannel(wsconn *websocket.Conn) *Channel {
	ch := &Channel{
		id:     -1,
		conn:   pkg.NewConnWs(wsconn, c.pool),
		recvCh: make(chan *msg.Msg, 100),
		sendCh: make(chan *msg.Msg, 100),
	}
	// ch.serveIO()
	return ch
}

// func (c *Channel) serveIO() {
// 	go c.recvRoutine()
// 	go c.sendRoutine()
// }

// func (c *Channel) recvRoutine() {
// 	for {
// 		message, err := c.conn.ReadMsg()
// 		if err != nil {
// 			// TODO
// 			continue
// 		}
// 		c.recvCh <- message
// 	}
// }

// func (c *Channel) sendRoutine() {
// 	for message := range c.sendCh {
// 		err := c.conn.WriteMsg(message)
// 		if err != nil {
// 			// TODO
// 			continue
// 		}
// 	}
// }

func (c *Channel) Close() {
	c.done = true
	c.conn.Close()
}
