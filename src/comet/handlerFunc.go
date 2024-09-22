package comet

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/dao/localCache"
	"laneIM/src/pkg/laneLog"
	"strconv"
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
		laneLog.Logger.Fatalln("[server] faild to get token")
		ch.ForceClose()
		return
	}
	rt, err := c.pickLogic().Client.Auth(context.Background(), &logic.AuthReq{
		Params:    authReq.Params,
		CometAddr: c.conf.WindowIP + c.conf.GrpcPort,
		Userid:    authReq.Userid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to auth logic")
		return
	}

	if !rt.Pass {
		laneLog.Logger.Infoln("[server] reject auth")
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
// 		laneLog.Logger.Fatalln("[server] faild to new user", err)
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
		laneLog.Logger.Fatalln("[server] faild to decode userid", err)
		return
	}
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: []int64{croomidReq.Userid},
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to query logic room", err)
		return
	}
	ch.id = croomidReq.Userid
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
		laneLog.Logger.Fatalln("[server] marchal err", err)
		return
	}
	ch.Reply(outData, m.Seq, m.Path)

}

type BatchStructQueryRoom struct {
	arg *msg.CRoomidReq
	ch  *Channel
	seq int64
}

func (c *Comet) doQueryRoomBatch(in []*BatchStructQueryRoom) {
	userids := make([]int64, len(in))
	for i := range in {
		userids[i] = in[i].arg.Userid
	}
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: userids,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to query logic room", err)
		return
	}
	for i := range in {
		outstruct := &msg.CRoomidResp{
			Roomid: rt.Roomids[i].Roomid,
		}
		if len(rt.Roomids) != 0 {
			for _, roomid := range rt.Roomids[i].Roomid {
				c.Bucket(roomid).PutChannel(roomid, in[i].ch)
				// laneLog.Logger.Infoln("userid:", ch.id, "in room", roomid)
			}
		}
		// comet初始化已加入的room
		outData, err := proto.Marshal(outstruct)
		if err != nil {
			laneLog.Logger.Fatalln("[server] marchal err", err)
			return
		}
		in[i].ch.Reply(outData, in[i].seq, "queryRoom")
	}

}

func (c *Comet) HandleQueryRoomBatch(m *msg.Msg, ch *Channel) {
	croomidReq := &msg.CRoomidReq{}
	err := proto.Unmarshal(m.Data, croomidReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode userid", err)
		return
	}
	ch.id = croomidReq.Userid
	c.BatcherQueryRoom.Add(&BatchStructQueryRoom{
		arg: croomidReq,
		ch:  ch,
		seq: m.Seq,
	})
}

// func (c *Comet) HandleSendRoom(m *msg.Msg, ch *Channel) {
// 	cSendRoomReq := &msg.CSendRoomReq{}
// 	err := proto.Unmarshal(m.Data, cSendRoomReq)
// 	if err != nil {
// 		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
// 		return
// 	}

// 	err = c.LogictSendMsg(&logic.SendMsgReq{
// 		Data:   []byte(cSendRoomReq.Msg),
// 		Path:   m.Path,
// 		Addr:   c.conf.WindowIP+c.conf.GrpcPort,
// 		Userid: cSendRoomReq.Userid,
// 		Roomid: cSendRoomReq.Roomid,
// 	})
// 	if err != nil {
// 		laneLog.Logger.Fatalln("[server] faild to send logic", err)
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
	// start := time.Now()
	msgs := make([]*msg.SendMsgReq, len(in))
	for i := range in {
		msgs[i] = &msg.SendMsgReq{
			Data:      []byte(in[i].arg.Msg),
			Path:      "sendRoom",
			Addr:      c.conf.WindowIP + c.conf.GrpcPort,
			Userid:    in[i].arg.Userid,
			Roomid:    in[i].arg.Roomid,
			Userseq:   in[i].seq,
			Timeunix:  in[i].timeunix,
			Messageid: c.msgUUIDGenerator.Generator(),
		}
	}

	err := c.LogictSendMsgBatch(&msg.SendMsgBatchReq{
		Msgs: msgs,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		// for i := range in {
		// 	in[i].ch.Reply([]byte("not ack"), in[i].seq, "sendRoom")
		// }
		// laneLog.Logger.Debugf("doSendRoomBatch not ack count = %d spand %v", len(in), time.Since(start))
		return
	}
	// for i := range in {
	// 	in[i].ch.Reply([]byte("ack"), in[i].seq, "sendRoom")
	// }
	// laneLog.Logger.Debugf("doSendRoomBatch count = %d spand %v", len(in), time.Since(start))
}

func (c *Comet) HandleSendRoomBatch(m *msg.Msg, ch *Channel) {
	cSendRoomReq := &msg.CSendRoomReq{}
	err := proto.Unmarshal(m.Data, cSendRoomReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}
	c.BatcherSendRoom.Add(&BatchStructSendRoom{
		arg: cSendRoomReq,
		ch:  ch,
		seq: m.Seq,
	})
	ch.Reply([]byte("ack"), m.Seq, "sendRoom")
}

func (c *Comet) HandleNewUser(m *msg.Msg, ch *Channel) {
	cNewUserReq := &msg.CNewUserReq{}
	err := proto.Unmarshal(m.Data, cNewUserReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}

	rt, err := c.pickLogic().Client.NewUser(context.Background(), &logic.NewUserReq{})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}
	reply, err := proto.Marshal(&msg.CNewUserResp{
		Userid: rt.Userid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to encode proto")
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
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}

	for i := range rt.Userid {
		reply, err := proto.Marshal(&msg.CNewUserResp{
			Userid: rt.Userid[i],
		})
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to encode proto")
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
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}
	c.BatcherNewUser.Add(&BatchStructNewUser{arg: cNewUserReq, ch: ch, seq: m.Seq})
}

func (c *Comet) HandleNewRoom(m *msg.Msg, ch *Channel) {
	cNewRoomReq := &msg.CNewRoomReq{}
	err := proto.Unmarshal(m.Data, cNewRoomReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}

	rt, err := c.pickLogic().Client.NewRoom(context.Background(), &logic.NewRoomReq{
		Userid:    cNewRoomReq.Userid,
		CometAddr: c.conf.WindowIP + c.conf.GrpcPort,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}
	reply, err := proto.Marshal(&msg.CNewRoomResp{
		Roomid: rt.Roomid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to encode proto")
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
		Comet:  c.conf.WindowIP + c.conf.GrpcPort,
	})

	// laneLog.Logger.Debugln("join room spand:", time.Since(join))
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err, "addr", c.conf.WindowIP+c.conf.GrpcPort)
		return
	}
	for i, rid := range roomid {
		// laneLog.Logger.Debugln("put channel ch.id=", in[i].ch.id)
		c.Bucket(rid).PutChannel(rid, in[i].ch)
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
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
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
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.JoinRoom(context.Background(), &logic.JoinRoomReq{
		Userid: cJoinRoomReq.Userid,
		Roomid: cJoinRoomReq.Roomid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
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
		Server: c.conf.WindowIP + c.conf.GrpcPort,
	})
	// laneLog.Logger.Debugln("setonline batch spand:", time.Since(start))
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}
	query := time.Now()
	// 缓存中不存在的才查询
	mayExist, localnonexist := localCache.UserRoomidBatch(c.cache, userid)
	var needQueryUserid []int64
	for i, e := range localnonexist {
		if e {
			needQueryUserid = append(needQueryUserid, userid[i])
		}
	}

	laneLog.Logger.Warnf("userroom local cache 命中[%d/%d]", len(userid)-len(needQueryUserid), len(userid))
	var rt *logic.QueryRoomResp
	if len(needQueryUserid) != 0 {
		// laneLog.Logger.Warnf("查询远端 userid[%v] ", userid)
		rt, err = c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
			Userid: needQueryUserid,
		})
		laneLog.Logger.Debugln("doQueryRoom spand:", time.Since(query))
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to query logic room", err)
			return
		}
	}

	var retdata []byte

	// putchannel
	var index = 0
	for i := range mayExist {
		if localnonexist[i] {
			localCache.SetUserRoomid(c.cache, userid[i], rt.Roomids[index].Roomid)
			// laneLog.Logger.Warnf("远端获取 userid[%d]  roomid[%v]", userid[i], rt.Roomids[index])
			for _, roomid := range rt.Roomids[index].Roomid {
				// laneLog.Logger.Debugln("put channel ch.id=", in[i].ch.id)
				c.Bucket(roomid).PutChannel(roomid, in[i].ch)
			}
			retdata, err = proto.Marshal(&msg.COnlineResp{
				Ack:    true,
				Roomid: rt.Roomids[index].Roomid,
			})
			index++
		} else {
			// laneLog.Logger.Warnf("localcache userid[%d]  roomid[%v]", userid[i], mayExist[i])
			for _, roomid := range mayExist[i] {
				// laneLog.Logger.Debugln("put channel ch.id=", in[i].ch.id)
				c.Bucket(roomid).PutChannel(roomid, in[i].ch)
			}
			retdata, err = proto.Marshal(&msg.COnlineResp{
				Ack:    true,
				Roomid: mayExist[i],
			})
		}

		if err != nil {
			laneLog.Logger.Errorln("faild marshal proto", err)
		}
		in[i].ch.Reply([]byte(retdata), in[i].seq, "online")
	}
	laneLog.Logger.Debugln("doSetOnlineBatch spand:", time.Since(start))
}

func (c *Comet) HandleSetOnlineBatch(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COnlineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}
	ch.id = COnlineReq.Userid
	// laneLog.Logger.Debugln("ch.id=", ch.id, "COnlineReq.Userid=", COnlineReq.Userid)
	c.BatcherSetOnline.Add(&BatchStructSetOnline{arg: COnlineReq, seq: m.Seq, ch: ch})
}

func (c *Comet) HandleSetOnline(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COnlineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.SetOnline(context.Background(), &logic.SetOnlineReq{
		Userid: COnlineReq.Userid,
		Server: c.conf.WindowIP + c.conf.GrpcPort,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}

	ch.id = COnlineReq.Userid
	{ // putchannel

		rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
			Userid: []int64{COnlineReq.Userid},
		})
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to query logic room", err)
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
		Server: c.conf.WindowIP + c.conf.GrpcPort,
	})
	// laneLog.Logger.Debugln("setonline batch spand:", time.Since(start))
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}

	// del channel
	c.DelChannelBatch(in)

	laneLog.Logger.Debugln("doSetOfflineBatch spand:", time.Since(start))
}
func (c *Comet) HandleSetOfflineBatch(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COfflineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faildto decode proto", err)
		return
	}
	ch.id = COnlineReq.Userid
	c.BatcherSetOffline.Add(&BatchStructSetOffline{arg: COnlineReq, seq: m.Seq, ch: ch})
}

func (c *Comet) HandleSetOffline(m *msg.Msg, ch *Channel) {
	COnlineReq := &msg.COfflineReq{}
	err := proto.Unmarshal(m.Data, COnlineReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faildto decode proto", err)
		return
	}

	_, err = c.pickLogic().Client.SetOffline(context.Background(), &logic.SetOfflineReq{
		Userid: COnlineReq.Userid,
		Server: c.conf.WindowIP + c.conf.GrpcPort,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faildto send logic", err)
		return
	}

	{ // del channel
		c.DelChannel(ch)
	}

	// ch.Reply([]byte("ack"), m.Seq, m.Path)
}

type ChAndSeq struct {
	ch  *Channel
	seq int64
}

type BatchStructQueryStoreMsg struct {
	arg      *msg.CQueryStoreMessageReq
	chAndSeq *ChAndSeq
}

func (c *Comet) doQueryStoreMsgBatch(in []*BatchStructQueryStoreMsg) {
	// start := time.Now()
	multiPageInfos := make([]*msg.QueryMultiRoomPagesReq_RoomMultiPageInfo, 0)

	// roomid,messageid,size -> ch
	hashToCh := make(map[string]*ChAndSeq)
	// 聚合同一个room的pageinfos
	roomInfosMap := make(map[int64][]*msg.QueryMultiRoomPagesReq_RoomMultiPageInfo_PageInfo)
	for _, a := range in {
		roomInfosMap[a.arg.Roomid] = append(roomInfosMap[a.arg.Roomid], &msg.QueryMultiRoomPagesReq_RoomMultiPageInfo_PageInfo{
			MessageId: a.arg.MessageId,
			Size:      a.arg.Size,
		})
		// laneLog.Logger.Debugln("hash:", HashRoomidMessageidSize(a.arg.Roomid, a.arg.MessageId, int(a.arg.Roomid)), "roomid=", a.arg.Roomid, "msgid=", a.arg.MessageId, "")
		hashToCh[HashRoomidMessageidSize(a.arg.Roomid, a.arg.MessageId, int(a.arg.Size))] = a.chAndSeq
	}
	for roomid, pageinfo := range roomInfosMap {
		multiPageInfos = append(multiPageInfos, &msg.QueryMultiRoomPagesReq_RoomMultiPageInfo{
			Roomid:    roomid,
			PageInfos: pageinfo,
		})
	}

	rt, err := c.pickLogic().Client.QueryStoreMsgBatch(context.Background(), &msg.QueryMultiRoomPagesReq{
		RoomMultiPageInfos: multiPageInfos,
	})
	for i, pages := range rt.RoomMultiPageMsgs {
		roomid := multiPageInfos[i].Roomid
		pageInfos := multiPageInfos[i].PageInfos
		pageMsgs := pages.PagesMsgs
		for i, page := range pageInfos {
			CH := hashToCh[HashRoomidMessageidSize(roomid, page.MessageId, int(page.Size))]
			data, err := proto.Marshal(pageMsgs[i])
			if err != nil {
				laneLog.Logger.Fatalln(err)
			}
			CH.ch.Reply(data, CH.ch.id, "his")
		}
	}
	// laneLog.Logger.Debugln("setonline batch spand time :", time.Since(start))

	// 再次分发给不同的user

	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send logic", err)
		return
	}
}

func (c *Comet) HandleHisMsg(m *msg.Msg, ch *Channel) {
	CHisMsgReq := &msg.CQueryStoreMessageReq{}
	err := proto.Unmarshal(m.Data, CHisMsgReq)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faildto decode proto", err)
		return
	}
	c.BatcherHis.Add(&BatchStructQueryStoreMsg{arg: CHisMsgReq, chAndSeq: &ChAndSeq{
		ch:  ch,
		seq: m.Seq,
	}})
}

func HashRoomidMessageidSize(roomid, messageid int64, size int) string {
	// 将参数转换为字符串并拼接
	data := strconv.FormatInt(roomid, 36) + strconv.FormatInt(messageid, 36) + strconv.Itoa(size)

	// 创建 SHA-256 哈希对象
	hash := sha256.New()
	// 写入数据
	hash.Write([]byte(data))
	// 计算哈希值
	hashBytes := hash.Sum(nil)

	// 返回十六进制编码的哈希值
	return hex.EncodeToString(hashBytes)
}

type BatcheStructLast struct {
	roomid int64
	ch     *Channel
	seq    int64
}

func (c *Comet) doQueryLast(in []*BatcheStructLast) {
	roomids := make([]int64, len(in))
	for i := range in {
		roomids[i] = in[i].roomid
	}
	rt, err := c.pickLogic().Client.QueryLast(context.Background(), &logic.QueryLastReq{
		Roomid: roomids,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faildto decode proto", err)
		return
	}
	for i := range rt.MessageId {
		in[i].ch.Reply([]byte(strconv.FormatInt(rt.MessageId[i], 10)), in[i].seq, "last")
	}
}

func (c *Comet) HandleLast(m *msg.Msg, ch *Channel) {
	roomid, _ := strconv.ParseInt(string(m.Data), 10, 64)
	c.BatcherLast.Add(&BatcheStructLast{
		ch:     ch,
		roomid: roomid,
		seq:    m.Seq,
	})
}

func (c *Comet) HandleSubscribe(m *msg.Msg, ch *Channel) {
	// laneLog.Logger.Debugln("pass")
	roomid, err := strconv.ParseInt(string(m.Data), 10, 64)
	if err != nil {
		laneLog.Logger.Fatalln(err)
	}
	c.Bucket(roomid).PutFullChannel(roomid, ch)
	ch.Reply([]byte("ack"), m.Seq, "subon")
}

func (c *Comet) HandleSubscribeOff(m *msg.Msg, ch *Channel) {
	roomid, err := strconv.ParseInt(string(m.Data), 10, 64)
	if err != nil {
		laneLog.Logger.Fatalln(err)
	}
	c.Bucket(roomid).DelFullChannel(roomid, ch)
	ch.Reply([]byte("ack"), m.Seq, "suboff")
}
