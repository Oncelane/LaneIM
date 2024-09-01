package comet

import (
	"context"
	"encoding/json"
	"laneIM/proto/logic"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 协议升级
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Token struct {
	Content string `json:"token"`
}

func (c *Comet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 建立websocket链接
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upGrader fail", err)
		return
	}
	// 通过解析网址，加入选定群组
	ch := c.NewChannel(ws)

	// receive token
	message, err := ch.conn.ReadMsg()
	if err != nil {
		log.Println("faild to get token")
		w.Write([]byte("faild to get token"))
		ch.conn.Close()
		return
	}
	token := &Token{}
	err = json.Unmarshal(message.Byte, token)
	if err != nil {
		log.Println("faild to decode token")
		w.Write([]byte("faild to decode token"))
		ch.conn.Close()
		return
	}

	// send token to logic
	authReq := &logic.AuthReq{
		Token: []byte("this is token message password:XXX userid:XXX"),
	}

	rt, err := c.pickLogic().Client.Auth(context.Background(), authReq)
	if err != nil {
		log.Println("faild to auth logic")
		w.Write([]byte("faild to auth logic"))
		ch.conn.Close()
		return
	}

	if !rt.Pass {
		log.Println("reject auth")
		w.Write([]byte("faild to auth logic"))
		ch.conn.Close()
		return
	}

	// success auth
	ch.id = rt.Userid
	log.Println("user id:", rt.Userid, "auth success")

	// 查询 userid  所处 room

	roomReq := &logic.QueryRoomReq{
		Userid: []int64{ch.id},
	}
	roomResp, err := c.pickLogic().Client.QueryRoom(context.Background(), roomReq)
	if err != nil {
		log.Println("faild to query userid", ch.id, "'s roomid")
		return
	}
	// 加入room
	if len(roomResp.Roomids) != 0 {
		for _, roomid := range roomResp.Roomids[0].Roomid {
			c.Bucket(ch.id).PutChannel(roomid, ch)
			log.Println("userid:", ch.id, "in room", roomid)
		}
	}

	return

}
