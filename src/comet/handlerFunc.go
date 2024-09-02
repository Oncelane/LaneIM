package comet

import (
	"context"
	"encoding/json"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"log"
)

type UserJson struct {
	Userid int64
}

func (c *Comet) TestNewUser(m *msg.Msg) {
	in := &logic.NewUserReq{}
	rt, err := c.pickLogic().Client.NewUser(context.Background(), in)
	if err != nil {
		log.Println("faild to new user", err)
		return
	}
	newUserJson := UserJson{Userid: rt.Userid}
	outdata, err := json.Marshal(newUserJson)
	out := &msg.Msg{
		Path: m.Path,
		Seq:  m.Seq,
		Byte: outdata,
	}
	c.chmu.RLock()
	err = c.channels[rt.Userid].conn.WriteMsg(out)
	if err != err {
		log.Println("faild to response Newuser to webscoket", err)
	}
	c.chmu.RUnlock()
}

func (c *Comet) TestRoom(m *msg.Msg) {
	useridStruct := &UserJson{}
	err := json.Unmarshal(m.Byte, useridStruct)
	if err != nil {
		log.Println("faild to decode userid", err)
	}

	in := &logic.QueryRoomReq{
		Userid: []int64{useridStruct.Userid},
	}
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), in)
	if err != nil {
		log.Println("faild to query logic room", err)
	}
	out := &msg.Msg{
		Path: m.Path,
		Seq:  m.Seq,
		Byte: outdata,
	}

	c.chmu.RLock()
	c.channels[useridStruct.Userid].conn.WriteMsg()
	c.chmu.RUnlock()

}
