package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"

	"github.com/gorilla/websocket"
)

type Channel struct {
	id   int64
	conn pkg.MsgReadWriteCloser
	// conn   pkg.MsgReadWriteCloser
	sendCh chan *msg.MsgBatch
	done   bool
}

func (c *Comet) NewChannel(wsconn *websocket.Conn) *Channel {
	ch := &Channel{
		id:     -1,
		conn:   pkg.NewConnWs(wsconn, c.pool),
		sendCh: make(chan *msg.MsgBatch, 100),
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
	// laneLog.Logger.Debugln("ch.id%d start receive", ch.id)
	for {
		message, err := ch.conn.ReadMsg()
		if ch.done {
			return
		}
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				laneLog.Logger.Fatalln("[websocket] faild to get ws message")
				c.DelChannel(ch)
				return
			}
			// laneLog.Logger.Infoln("websocket close", ch.id)
			c.DelChannel(ch)
			return
		}
		for i := range message.Msgs {
			// laneLog.Logger.Infoln("message.Path", message.Msgs[i].Path)
			f := c.funcRout.Find(message.Msgs[i].Path)
			if f == nil {
				laneLog.Logger.Fatalln("[websocket] wrong method")
				continue
			}
			go f(message.Msgs[i], ch)
			if ch.done {
				return
			}
		}
	}
}

func (c *Comet) sendRoutine(ch *Channel) {
	for message := range ch.sendCh {
		// laneLog.Logger.Infoln("send room msg to websocket")
		err := ch.conn.WriteMsg(message)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				laneLog.Logger.Fatalln("[server] faild to get ws message", err)
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
	c.sendCh <- &msg.MsgBatch{
		Msgs: []*msg.Msg{
			{
				Data: data,
				Seq:  seq,
				Path: path,
			},
		},
	}
}

func (c *Channel) Send(in *msg.MsgBatch) {
	c.sendCh <- in
}

func (c *Channel) Close() {
	c.done = true
	// close(c.sendCh)
	c.conn.Close()
}

func (c *Channel) ForceClose() {
	c.done = true
	// close(c.sendCh)
	c.conn.ForceClose()
}
