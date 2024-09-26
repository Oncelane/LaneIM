package job

import (
	"laneIM/proto/msg"
	"sort"
)

type Timeunix []*msg.SendMsgReq

func (a Timeunix) Len() int           { return len(a) }
func (a Timeunix) Less(i, j int) bool { return a[i].Timeunix.AsTime().Before(a[j].Timeunix.AsTime()) }
func (a Timeunix) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (j *Job) SortRoomMsg_Timeunix(roomMsg *msg.SendMsgBatchReq) {
	sort.Sort(Timeunix(roomMsg.Msgs))
}
