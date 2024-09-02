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
	c.chmu.RLock()
	err = c.channels[rt.Userid].conn.WriteMsg(&msg.Msg{
		Path: m.Path,
		Seq:  m.Seq,
		Data: outdata,
	})
	if err != err {
		log.Println("faild to response Newuser to webscoket", err)
	}
	c.chmu.RUnlock()
}

func (c *Comet) TestRoom(m *msg.Msg) {
	useridStruct := &UserJson{}
	err := json.Unmarshal(m.Data, useridStruct)
	if err != nil {
		log.Println("faild to decode userid", err)
		return
	}

	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: []int64{useridStruct.Userid},
	})
	if err != nil {
		log.Println("faild to query logic room", err)
		return
	}
	outstruct := struct {
		Roomid []int64
	}{
		Roomid: rt.Roomids[0].Roomid,
	}
	outData, err := json.Marshal(outstruct)
	c.chmu.RLock()
	c.channels[useridStruct.Userid].conn.WriteMsg(&msg.Msg{
		Path: m.Path,
		Seq:  m.Seq,
		Data: outData,
	})
	c.chmu.RUnlock()

}

func (c *Comet) TestSendRoom(m *msg.Msg) {
	MsgRoomid := &struct {
		Msg    string
		Roomid int64
		Userid int64
	}{}
	err := json.Unmarshal(m.Data, MsgRoomid)
	if err != nil {
		log.Println("faild to decode json", err)
		return
	}

	_, err = c.pickLogic().Client.SendMsg(context.Background(), &logic.SendMsgReq{
		Data:   []byte(MsgRoomid.Msg),
		Path:   "room",
		Addr:   c.conf.Addr,
		Userid: MsgRoomid.Userid,
		Roomid: MsgRoomid.Userid,
	})
	if err != nil {
		log.Println("faild to send logic", err)
		return
	}
	c.chmu.RLock()
	c.channels[MsgRoomid.Userid].conn.WriteMsg(&msg.Msg{
		Path: m.Path,
		Seq:  m.Seq,
		Data: []byte("ack"),
	})
	c.chmu.RUnlock()
}
