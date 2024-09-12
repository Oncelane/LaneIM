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
				return &msg.MsgBatch{}
			},
		},
	}
}

func (m *MsgPool) Put(in *msg.MsgBatch) {
	in.Reset()
	m.pool.Put(in)
}

func (m *MsgPool) Get() *msg.MsgBatch {
	return m.pool.Get().(*msg.MsgBatch)
}
