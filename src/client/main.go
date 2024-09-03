package main

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn    *pkg.ConnWs
	userid  int64
	roomids []int64
	online  bool
	message string
}

var cometAddr string = "ws://localhost:40051/ws"
var clients []*Client

func main() {
	clients = make([]*Client, 4)
	for i := range clients {
		clients[i] = NewClient(21 + int64(i))
		clients[i].Connect(cometAddr)
		clients[i].Auth("i am 2" + strconv.FormatInt(int64(i+1), 10))
		clients[i].message = "hello i am 2" + strconv.FormatInt(int64(i+1), 10)
		clients[i].Receive()
	}
	select {}
}

func NewClient(userid int64) *Client {
	return &Client{
		userid: userid,
		online: false,
	}
}

func (c *Client) Connect(cometAddr string) {
	// 连接到WebSocket服务
	conn, _, err := websocket.DefaultDialer.Dial(cometAddr, nil)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	log.Println("连接到comet:", cometAddr)
	c.conn = pkg.NewConnWs(conn, pkg.NewMsgPool())
}

func (c *Client) Auth(str string) {
	// 发送消息到服务器
	token := &msg.CAuthReq{
		Params: []string{str},
		Userid: c.userid,
	}
	tokenDate, err := proto.Marshal(token)
	if err != nil {
		log.Println("faild to encode token")
	}
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "auth",
		Seq:  1,
		Data: tokenDate,
	})
	if err != nil {
		log.Println("发送消息错误:", err)
	}
}

func (c *Client) SendRoomMsg(message string) {
	sendRoomReq := &msg.CSendRoomReq{
		Userid: c.userid,
		Roomid: c.roomids[0],
		Msg:    message,
	}
	data, err := proto.Marshal(sendRoomReq)
	if err != nil {
		log.Println("faild to proto marshal", err)
	}
	log.Println("发送房间消息")
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "sendRoom",
		Seq:  1,
		Data: data,
	})
	if err != nil {
		log.Println("send err", err)
	}
}

func (c *Client) QueryRoom() {
	cRoomidReq := &msg.CRoomidReq{
		Userid: c.userid,
	}
	data, err := proto.Marshal(cRoomidReq)
	if err != nil {
		log.Println("faild to marhal")
		return
	}

	err = c.conn.WriteMsg(&msg.Msg{
		Path: "queryRoom",
		Seq:  1,
		Data: data,
	})
	if err != nil {
		log.Println("send err:", err)
	}
}

func (c *Client) Receive() {
	go func() {
		for {
			message, err := c.conn.ReadMsg()
			if err != nil {
				log.Println("comet error:", err)
				return
			}
			// log.Printf("comet reply: %s", message.String())
			switch message.Path {
			case "auth":
				log.Println("auth:", string(message.Data))
				c.QueryRoom()
			case "sendRoom":

				log.Println("send room:", string(message.Data))
			case "queryRoom":
				roomResp := &msg.CRoomidResp{}
				err := proto.Unmarshal(message.Data, roomResp)
				if err != nil {
					log.Println("wrong", message.Path, err)
					continue
				}
				c.roomids = roomResp.Roomid
				c.SendRoomMsg(c.message)
				log.Println("query room:", c.roomids[0])
			case "roomMsg":
				log.Printf("ch.id[%d] roomMsg receive:%s\n", c.userid, string(message.Data))
			}
		}
	}()
}
