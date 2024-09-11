package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog.go"

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
		done:   false,
	}
	// ch.serveIO()
	return ch
}

func (c *Comet) serveIO(ch *Channel) {
	go c.recvRoutine(ch)
	go c.sendRoutine(ch)
}

func (c *Comet) recvRoutine(ch *Channel) {
	for {
		message, err := ch.conn.ReadMsg()
		if ch.done {
			return
		}
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				laneLog.Logger.Infoln("faild to get ws message")
				c.DelChannel(ch)
				return
			}
			// laneLog.Logger.Infoln("websocket close", ch.id)
			c.DelChannel(ch)
			return
		}
		// laneLog.Logger.Infoln("message.Path", message.Path)
		f := c.funcRout.Find(message.Path)
		if f == nil {
			laneLog.Logger.Infoln("wrong method")
			continue
		}
		go f(message, ch)
		if ch.done {
			return
		}
	}
}

func (c *Comet) sendRoutine(ch *Channel) {
	for message := range ch.sendCh {
		// laneLog.Logger.Infoln("send room msg to websocket")
		err := ch.conn.WriteMsg(message)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				laneLog.Logger.Infoln("faild to get ws message")
				c.DelChannel(ch)
				return
			}
			// laneLog.Logger.Infoln("websocket close", ch.id)
			c.DelChannel(ch)
			return
		}
	}
}

func (c *Channel) Reply(data []byte, seq int64, path string) {
	c.sendCh <- &msg.Msg{
		Data: data,
		Seq:  seq,
		Path: path,
	}
}

func (c *Channel) Close() {
	c.done = true
	// close(c.sendCh)
	c.conn.Close()
}
