package comet

import (
	"context"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/pkg/laneLog.go"
	"time"

	"google.golang.org/protobuf/proto"
)

type UserJson struct {
	Userid int64
}

func (c *Comet) HandleAuth(m *msg.Msg, ch *Channel) {
	authReq := &msg.CAuthReq{}
	err := proto.Unmarshal(m.Data, authReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to get token")
		ch.conn.Close()
		return
	}
	rt, err := c.pickLogic().Client.Auth(context.Background(), &logic.AuthReq{
		Params:    authReq.Params,
		CometAddr: c.conf.Addr,
		Userid:    authReq.Userid,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to auth logic")
		return
	}

	if !rt.Pass {
		laneLog.Logger.Infoln("reject auth")
		ch.Reply([]byte("false"), m.Seq, m.Path)
		return
	}

	// success auth
	// laneLog.Logger.Infoln("user id:", authReq.Userid, "auth success")
	ch.id = authReq.Userid
	ch.Reply([]byte("true"), m.Seq, m.Path)
}

// func (c *Comet) HandleNewUser(m *msg.Msg, ch *Channel) {
// 	in := &logic.NewUserReq{}
// 	rt, err := c.pickLogic().Client.NewUser(context.Background(), in)
// 	if err != nil {
// 		laneLog.Logger.Infoln("faild to new user", err)
// 		return
// 	}
// 	newUserJson := UserJson{Userid: rt.Userid}
// 	outdata, err := json.Marshal(newUserJson)
// 	ch.Reply(outdata, m.Seq, m.Path)
// }

func (c *Comet) HandleRoom(m *msg.Msg, ch *Channel) {
	croomidReq := &msg.CRoomidReq{}
	err := proto.Unmarshal(m.Data, croomidReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode userid", err)
		return
	}
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: []int64{croomidReq.Userid},
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to query logic room", err)
		return
	}
	if len(rt.Roomids) != 0 {
		for _, roomid := range rt.Roomids[0].Roomid {
			c.Bucket(roomid).PutChannel(roomid, ch)
			// laneLog.Logger.Infoln("userid:", ch.id, "in room", roomid)
		}
	}
	outstruct := &msg.CRoomidResp{
		Roomid: rt.Roomids[0].Roomid,
	}

	// comet初始化已加入的room

	outData, err := proto.Marshal(outstruct)
	if err != nil {
		laneLog.Logger.Infoln("marchal err", err)
		return
	}
	ch.Reply(outData, m.Seq, m.Path)

}

func (c *Comet) HandleSendRoom(m *msg.Msg, ch *Channel) {
	cSendRoomReq := &msg.CSendRoomReq{}
	err := proto.Unmarshal(m.Data, cSendRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	err = c.LogictSendMsg(&logic.SendMsgReq{
		Data:   []byte(cSendRoomReq.Msg),
		Path:   m.Path,
		Addr:   c.conf.Addr,
		Userid: cSendRoomReq.Userid,
		Roomid: cSendRoomReq.Roomid,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	ch.Reply([]byte("ack"), m.Seq, m.Path)
}

func (c *Comet) HandleNewUser(m *msg.Msg, ch *Channel) {
	cNewUserReq := &msg.CNewUserReq{}
	err := proto.Unmarshal(m.Data, cNewUserReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	rt, err := c.pickLogic().Client.NewUser(context.Background(), &logic.NewUserReq{})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	reply, err := proto.Marshal(&msg.CNewUserResp{
		Userid: rt.Userid,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to encode proto")
		return
	}
	ch.id = rt.Userid
	ch.Reply(reply, m.Seq, m.Path)
}

type BatchStructNewUser struct {
	arg *msg.CNewUserReq
	seq int64
	ch  *Channel
}

func (c *Comet) doNewUserBatch(in []*BatchStructNewUser) {
	rt, err := c.pickLogic().Client.NewUserBatch(context.Background(), &logic.NewUserBatchReq{Count: int64(len(in))})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}

	for i := range rt.Userid {
		reply, err := proto.Marshal(&msg.CNewUserResp{
			Userid: rt.Userid[i],
		})
		if err != nil {
			laneLog.Logger.Infoln("faild to encode proto")
			return
		}
		in[i].ch.id = rt.Userid[i]
		in[i].ch.Reply(reply, in[i].seq, "newUser")
	}
	laneLog.Logger.Infoln("debug: doNewUserBatch run once")
}

func (c *Comet) HandleNewUserBatch(m *msg.Msg, ch *Channel) {
	cNewUserReq := &msg.CNewUserReq{}
	err := proto.Unmarshal(m.Data, cNewUserReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	c.BatcherNewUser.Add(&BatchStructNewUser{arg: cNewUserReq, ch: ch, seq: m.Seq})
}

func (c *Comet) HandleNewRoom(m *msg.Msg, ch *Channel) {
	cNewRoomReq := &msg.CNewRoomReq{}
	err := proto.Unmarshal(m.Data, cNewRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	rt, err := c.pickLogic().Client.NewRoom(context.Background(), &logic.NewRoomReq{
		Userid:    cNewRoomReq.Userid,
		CometAddr: c.conf.Addr,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	reply, err := proto.Marshal(&msg.CNewRoomResp{
		Roomid: rt.Roomid,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to encode proto")
		return
	}
	ch.Reply(reply, m.Seq, m.Path)
}

type BatchStructJoinRoom struct {
	arg *msg.CJoinRoomReq
	ch  *Channel
	seq int64
}

func (c *Comet) doJoinRoomBatch(in []*BatchStructJoinRoom) {
	userid := make([]int64, len(in))
	roomid := make([]int64, len(in))
	for i := range in {
		userid[i] = in[i].arg.Userid
		roomid[i] = in[i].arg.Roomid
	}
	join := time.Now()
	_, err := c.pickLogic().Client.JoinRoomBatch(context.Background(), &logic.JoinRoomBatchReq{
		Userid: userid,
		Roomid: roomid,
	})
	laneLog.Logger.Debugln("join room spand:", time.Since(join))
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	for i := range in {
		in[i].ch.Reply([]byte("ack"), in[i].seq, "joinRoom")
	}

}
func (c *Comet) HandleJoinRoomBatch(m *msg.Msg, ch *Channel) {
	cJoinRoomReq := &msg.CJoinRoomReq{}
	err := proto.Unmarshal(m.Data, cJoinRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	c.BatcherJoinRoom.Add(&BatchStructJoinRoom{
		arg: cJoinRoomReq,
		ch:  ch,
		seq: m.Seq,
	})

}

func (c *Comet) HandleJoinRoom(m *msg.Msg, ch *Channel) {
	cJoinRoomReq := &msg.CJoinRoomReq{}
	err := proto.Unmarshal(m.Data, cJoinRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.JoinRoom(context.Background(), &logic.JoinRoomReq{
		Userid: cJoinRoomReq.Userid,
		Roomid: cJoinRoomReq.Roomid,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	ch.Reply([]byte("ack"), m.Seq, m.Path)
}

type BatchStructSetOnline struct {
	arg *msg.COnlineReq
	seq int64
	ch  *Channel
}

func (c *Comet) doSetOnlineBatch(in []*BatchStructSetOnline) {
	start := time.Now()
	userid := make([]int64, len(in))
	for i := range len(in) {
		userid[i] = in[i].arg.Userid
	}

	_, err := c.pickLogic().Client.SetOnlineBatch(context.Background(), &logic.SetOnlineBatchReq{
		Userid: userid,
		Server: c.conf.Addr,
	})
	laneLog.Logger.Debugln("setonline batch spand:", time.Since(start))
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	query := time.Now()
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: userid,
	})
	laneLog.Logger.Debugln("doQueryRoom spand:", time.Since(query))
	if err != nil {
		laneLog.Logger.Infoln("faild to query logic room", err)
		return
	}

	var retdata []byte

	if len(rt.Roomids) != len(in) {
		for i := range in {
			laneLog.Logger.Errorf("faild some query room, ask %s but return %s", len(in), len(rt.Roomids))
			retdata, err = proto.Marshal(&msg.COnlineResp{
				Ack:    false,
				Roomid: nil,
			})
			if err != nil {
				laneLog.Logger.Errorln("faild marshal proto", err)
			}
			in[i].ch.Reply(retdata, in[i].seq, "online")
		}
		return
	}

	// putchannel
	for i := range in {

		for _, roomid := range rt.Roomids[i].Roomid {
			c.Bucket(roomid).PutChannel(roomid, in[i].ch)
		}
		retdata, err = proto.Marshal(&msg.COnlineResp{
			Ack:    true,
			Roomid: rt.Roomids[i].Roomid,
		})
		if err != nil {
			laneLog.Logger.Errorln("faild marshal proto", err)
		}
		in[i].ch.Reply([]byte(retdata), in[i].seq, "online")
	}
	laneLog.Logger.Debugln("total spand:", time.Since(start))
}

func (c *Comet) HandleSetOnlineBatch(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COnlineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	ch.id = COnlineReq.Userid
	c.BatcherSetOline.Add(&BatchStructSetOnline{arg: COnlineReq, seq: m.Seq, ch: ch})
}

func (c *Comet) HandleSetOnline(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COnlineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.SetOnline(context.Background(), &logic.SetOnlineReq{
		Userid: COnlineReq.Userid,
		Server: c.conf.Addr,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	// reply, err := proto.Marshal(&msg.CJoinRoomResp{
	// 	Ack: true,
	// })
	// if err != nil {
	// 	laneLog.Logger.Infoln("faild to encode proto")
	// 	return
	// }

	{ // putchannel
		ch.id = COnlineReq.Userid

		rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
			Userid: []int64{COnlineReq.Userid},
		})
		if err != nil {
			laneLog.Logger.Infoln("faild to query logic room", err)
			return
		}
		if len(rt.Roomids) != 0 {
			for _, roomid := range rt.Roomids[0].Roomid {
				c.Bucket(roomid).PutChannel(roomid, ch)
				// laneLog.Logger.Infoln("userid:", ch.id, "in room", roomid)
			}
		}
	}

	ch.Reply([]byte("ack"), m.Seq, m.Path)
}
