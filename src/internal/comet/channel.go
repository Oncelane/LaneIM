package comet

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"sync"

	"github.com/gorilla/websocket"
)

func (c *Comet) NewChannel(wsconn *websocket.Conn) *Channel {
	ch := c.poolCh.Get()
	ch.id = -1
	ch.done = false
	ch.conn = pkg.NewConnWs(wsconn, c.pool)
	ch.sendCh = make(chan *msg.MsgBatch, 2)
	ch.onceCh = sync.Once{}
	return ch
}

func (c *Comet) serveIO(ch *Channel) {
	go c.recvRoutine(ch)
	go c.sendRoutine(ch)
}

func (c *Comet) recvRoutine(ch *Channel) {
	// laneLog.Logger.Debugln("ch.id%d start receive", ch.id)
	// defer laneLog.Logger.Warnln("Read Shoud Be Quit!!!")
	defer c.DelChannel(ch)
	for {
		message, err := ch.conn.ReadMsg()
		if err != nil {
			// if _, ok := err.(*websocket.CloseError); !ok {
			// 	return
			// }
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
		}
	}

}

func (c *Comet) sendRoutine(ch *Channel) {
	// defer laneLog.Logger.Warnln("Send Shoud Be Quit!!!")
	defer c.DelChannel(ch)
	for message := range ch.sendCh {
		if ch.done {
			return
		}
		err := ch.conn.WriteMsg(message)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				// laneLog.Logger.Fatalln("[server] faild to get ws message", err)

				return
			}
			return
		}
	}
}

type Channel struct {
	id   int64
	conn pkg.MsgReadWriteCloser
	// conn   pkg.MsgReadWriteCloser
	sendCh chan *msg.MsgBatch
	onceCh sync.Once
	done   bool
}

type ChannelPool struct {
	pool sync.Pool
}

func NewChannelPool() *ChannelPool {
	return &ChannelPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Channel{}
			},
		},
	}
}

func (c *ChannelPool) Get() *Channel {
	return c.pool.Get().(*Channel)
}

func (c *ChannelPool) Put(in *Channel) {
	// in.conn = nil
	// in.sendCh = nil
	// in.onceCh = sync.Once{}
	c.pool.Put(in)
}

func (c *Channel) Reply(data []byte, seq int64, path string) {
	if c.done {
		return
	}
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
	if c.done {
		return
	}
	c.sendCh <- in
}

func (c *Channel) PassiveClose() error {
	return c.conn.PassiceClose()
}

func (c *Channel) ForceClose(p *ChannelPool) error {
	c.done = true
	var err error
	c.onceCh.Do(func() {
		// laneLog.Logger.Debugf("channel[%d] close send channel", c.id)
		close(c.sendCh)
		err = c.conn.ForceClose()
		p.Put(c)
	})
	return err
}
