package pkg

import (
	"laneIM/proto/msg"
	"sync"
)

type MsgPool struct {
	pool sync.Pool
}

func NewMsgPool() *MsgPool {
	return &MsgPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &msg.Msg{}
			},
		},
	}
}

func (m *MsgPool) Put(in *msg.Msg) {
	in.Reset()
	m.pool.Put(in)
}

func (m *MsgPool) Get() *msg.Msg {
	return m.pool.Get().(*msg.Msg)
}
