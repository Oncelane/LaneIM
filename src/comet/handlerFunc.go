package comet

import (
	"context"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"log"

	"google.golang.org/protobuf/proto"
)

type UserJson struct {
	Userid int64
}

func (c *Comet) HandleAuth(m *msg.Msg, ch *Channel) {
	authReq := &msg.CAuthReq{}
	err := proto.Unmarshal(m.Data, authReq)
	if err != nil {
		log.Println("faild to get token")
		ch.conn.Close()
		return
	}
	rt, err := c.pickLogic().Client.Auth(context.Background(), &logic.AuthReq{
		Params:    authReq.Params,
		CometAddr: c.conf.Addr,
		Userid:    authReq.Userid,
	})
	if err != nil {
		log.Println("faild to auth logic")
		return
	}

	if !rt.Pass {
		log.Println("reject auth")
		ch.Reply([]byte("false"), m.Seq, m.Path)
		return
	}

	// success auth
	log.Println("user id:", authReq.Userid, "auth success")
	ch.id = authReq.Userid
	ch.Reply([]byte("true"), m.Seq, m.Path)
}

// func (c *Comet) HandleNewUser(m *msg.Msg, ch *Channel) {
// 	in := &logic.NewUserReq{}
// 	rt, err := c.pickLogic().Client.NewUser(context.Background(), in)
// 	if err != nil {
// 		log.Println("faild to new user", err)
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
		log.Println("faild to decode userid", err)
		return
	}
	rt, err := c.pickLogic().Client.QueryRoom(context.Background(), &logic.QueryRoomReq{
		Userid: []int64{croomidReq.Userid},
	})
	if err != nil {
		log.Println("faild to query logic room", err)
		return
	}
	if len(rt.Roomids) != 0 {
		for _, roomid := range rt.Roomids[0].Roomid {
			c.Bucket(roomid).PutChannel(roomid, ch)
			log.Println("userid:", ch.id, "in room", roomid)
		}
	}
	outstruct := &msg.CRoomidResp{
		Roomid: rt.Roomids[0].Roomid,
	}

	// comet初始化已加入的room

	outData, err := proto.Marshal(outstruct)
	if err != nil {
		log.Println("marchal err", err)
		return
	}
	ch.Reply(outData, m.Seq, m.Path)

}

func (c *Comet) HandleSendRoom(m *msg.Msg, ch *Channel) {
	cSendRoomReq := &msg.CSendRoomReq{}
	err := proto.Unmarshal(m.Data, cSendRoomReq)
	if err != nil {
		log.Println("faild to decode json", err)
		return
	}

	_, err = c.pickLogic().Client.SendMsg(context.Background(), &logic.SendMsgReq{
		Data:   []byte(cSendRoomReq.Msg),
		Path:   m.Path,
		Addr:   c.conf.Addr,
		Userid: cSendRoomReq.Userid,
		Roomid: cSendRoomReq.Roomid,
	})
	if err != nil {
		log.Println("faild to send logic", err)
		return
	}
	ch.Reply([]byte("ack"), m.Seq, m.Path)
}
