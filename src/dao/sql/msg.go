package sql

import (
	"laneIM/proto/msg"
	"laneIM/src/model"
	"time"
)

func (d *SqlDB) SaveMessageBatch(in []*msg.SendMsgReq) {
	msgs := make([]model.RoomMessage, len(in))
	for i := range in {
		msgs[i].Content = string(in[i].Data)
		msgs[i].SenderId = in[i].Userid
		msgs[i].RoomId = in[i].Roomid
		msgs[i].SenderDate = time.Unix(in[i].Timeunix, 0)
	}
	d.DB.Save(&msgs)
}
