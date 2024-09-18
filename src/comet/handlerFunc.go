package comet

import (
	"context"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/dao/localCache"
	"laneIM/src/pkg/laneLog"
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

// func (c *Comet) HandleSendRoom(m *msg.Msg, ch *Channel) {
// 	cSendRoomReq := &msg.CSendRoomReq{}
// 	err := proto.Unmarshal(m.Data, cSendRoomReq)
// 	if err != nil {
// 		laneLog.Logger.Infoln("faild to decode proto", err)
// 		return
// 	}

// 	err = c.LogictSendMsg(&logic.SendMsgReq{
// 		Data:   []byte(cSendRoomReq.Msg),
// 		Path:   m.Path,
// 		Addr:   c.conf.Addr,
// 		Userid: cSendRoomReq.Userid,
// 		Roomid: cSendRoomReq.Roomid,
// 	})
// 	if err != nil {
// 		laneLog.Logger.Infoln("faild to send logic", err)
// 		return
// 	}
// 	ch.Reply([]byte("ack"), m.Seq, m.Path)
// }

type BatchStructSendRoom struct {
	arg      *msg.CSendRoomReq
	ch       *Channel
	seq      int64
	timeunix int64
}

func (c *Comet) doSendRoomBatch(in []*BatchStructSendRoom) {
	start := time.Now()
	msgs := make([]*msg.SendMsgReq, len(in))
	for i := range in {
		msgs[i] = &msg.SendMsgReq{
			Data:     []byte(in[i].arg.Msg),
			Path:     "sendRoom",
			Addr:     c.conf.Addr,
			Userid:   in[i].arg.Userid,
			Roomid:   in[i].arg.Roomid,
			Roomseq:  in[i].seq,
			Timeunix: in[i].timeunix,
		}
	}

	err := c.LogictSendMsgBatch(&msg.SendMsgBatchReq{
		Msgs: msgs,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	for i := range in {
		in[i].ch.Reply([]byte("ack"), in[i].seq, "sendRoom")
	}
	laneLog.Logger.Debugln("doSendRoomBatch spand", time.Since(start))
}

func (c *Comet) HandleSendRoomBatch(m *msg.Msg, ch *Channel) {
	cSendRoomReq := &msg.CSendRoomReq{}
	err := proto.Unmarshal(m.Data, cSendRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	c.BatcherSendRoom.Add(&BatchStructSendRoom{
		arg: cSendRoomReq,
		ch:  ch,
		seq: m.Seq,
	})

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
	start := time.Now()
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
	laneLog.Logger.Debugln("doNewUserBatch spand", time.Since(start))
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
	start := time.Now()
	userid := make([]int64, len(in))
	roomid := make([]int64, len(in))
	for i := range in {
		userid[i] = in[i].arg.Userid
		roomid[i] = in[i].arg.Roomid
		// laneLog.Logger.Warnf("send logic joinroom req : userid[%d] join room[%d]", userid[i], roomid[i])
	}
	// join := time.Now()
	_, err := c.pickLogic().Client.JoinRoomBatch(context.Background(), &logic.JoinRoomBatchReq{
		Userid: userid,
		Roomid: roomid,
	})

	// laneLog.Logger.Debugln("join room spand:", time.Since(join))
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	for i := range in {
		in[i].ch.Reply([]byte("ack"), in[i].seq, "joinRoom")
	}
	laneLog.Logger.Debugln("doJoinRoomBatch spand", time.Since(start))

}
func (c *Comet) HandleJoinRoomBatch(m *msg.Msg, ch *Channel) {
	cJoinRoomReq := &msg.CJoinRoomReq{}
	err := proto.Unmarshal(m.Data, cJoinRoomReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	// laneLog.Logger.Warnf("client msg req[%d] HandleJoinRoomBatch: userid[%d] join room[%d]", m.Seq, cJoinRoomReq.Userid, cJoinRoomReq.Roomid)
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
	// laneLog.Logger.Debugln("setonline batch spand:", time.Since(start))
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}
	// query := time.Now()
	// 缓存中不存在的才查询
	mayExist, localnonexist := localCache.UserRoomidBatch(c.cache, userid)
	var needQueryUserid []int64
	for i, e := range localnonexist {
		if e {
			needQueryUserid = append(needQueryUserid, userid[i])
		}
	}

	laneLog.Logger.Warnf("userroom local cache 命中[%d/%d]", len(userid)-len(needQueryUserid), len(userid))
	// laneLog.Logger.Warnf("查询远端 userid[%v] ", userid)
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: needQueryUserid,
	})
	// laneLog.Logger.Debugln("doQueryRoom spand:", time.Since(query))
	if err != nil {
		laneLog.Logger.Infoln("faild to query logic room", err)
		return
	}

	var retdata []byte

	// putchannel
	var index = 0
	for i := range mayExist {
		if localnonexist[i] {
			localCache.SetUserRoomid(c.cache, userid[i], rt.Roomids[index].Roomid)
			// laneLog.Logger.Warnf("远端获取 userid[%d]  roomid[%v]", userid[i], rt.Roomids[index])
			for _, roomid := range rt.Roomids[index].Roomid {
				c.Bucket(roomid).PutChannel(roomid, in[i].ch)
			}
			index++
		} else {
			// laneLog.Logger.Warnf("localcache userid[%d]  roomid[%v]", userid[i], mayExist[i])
			for _, roomid := range mayExist[i] {
				c.Bucket(roomid).PutChannel(roomid, in[i].ch)
			}
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
	c.BatcherSetOnline.Add(&BatchStructSetOnline{arg: COnlineReq, seq: m.Seq, ch: ch})
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

type BatchStructSetOffline struct {
	arg *msg.COfflineReq
	seq int64
	ch  *Channel
}

func (c *Comet) doSetOfflineBatch(in []*BatchStructSetOffline) {
	start := time.Now()
	userid := make([]int64, len(in))
	for i := range len(in) {
		userid[i] = in[i].arg.Userid
	}

	_, err := c.pickLogic().Client.SetOfflineBatch(context.Background(), &logic.SetOfflineBatchReq{
		Userid: userid,
		Server: c.conf.Addr,
	})
	// laneLog.Logger.Debugln("setonline batch spand:", time.Since(start))
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}

	// del channel

	// in[i].ch.Reply([]byte("ack"), in[i].seq, "offline")
	c.DelChannelBatch(in)

	laneLog.Logger.Debugln("total spand:", time.Since(start))
}
func (c *Comet) HandleSetOfflineBatch(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COfflineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}
	ch.id = COnlineReq.Userid
	c.BatcherSetOffline.Add(&BatchStructSetOffline{arg: COnlineReq, seq: m.Seq, ch: ch})
}

func (c *Comet) HandleSetOffline(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COfflineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.SetOffline(context.Background(), &logic.SetOfflineReq{
		Userid: COnlineReq.Userid,
		Server: c.conf.Addr,
	})
	if err != nil {
		laneLog.Logger.Infoln("faild to send logic", err)
		return
	}

	{ // del channel
		c.DelChannel(ch, false)
	}

	// ch.Reply([]byte("ack"), m.Seq, m.Path)
}
