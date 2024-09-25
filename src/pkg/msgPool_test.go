package pkg

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg/laneLog"
	"testing"
)

func TestMesagePool(t *testing.T) {
	p := NewMsgPool()
	ms := p.Get()

	ms.Msgs = make([]*msg.Msg, 3)

	for i := range ms.Msgs {
		ms.Msgs[i] = new(msg.Msg)
		ms.Msgs[i].Seq = int64(i)
		laneLog.Logger.Infoln(ms.Msgs[i].Seq)
	}

	tmp := ms.Msgs[0]
	laneLog.Logger.Infoln("ms[0]=", tmp.Seq)
	p.Put(ms)
	laneLog.Logger.Infoln(len(ms.Msgs))

	laneLog.Logger.Infoln("ms[0]=", tmp.Seq)

}
