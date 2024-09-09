package client

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Mgr          *ClientGroup
	conn         *pkg.ConnWs
	Userid       int64
	Roomids      []int64
	MsgCh        chan *string
	ReceiveCount int
}

// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
// var clients []*Client

// clients = make([]*Client, 4)
// for i := range clients {
// 	clients[i] = NewClient(21 + int64(i))
// 	clients[i].Connect(cometAddr[i%2])
// 	clients[i].Auth("i am 2" + strconv.FormatInt(int64(i+1), 10))
// 	clientQueryRoom.Add(1)
// 	go clients[i].Receive()
// }
// clientQueryRoom.Wait()
// for _, c := range clients {
// 	c.SendRoomMsg(c.message)
// }
// select {}

type ClientGroup struct {
	Num     int
	Clients []*Client
	Wait    sync.WaitGroup
}

func (c *ClientGroup) Send(msg *string) {
	for _, client := range c.Clients {
		//log.Println("in ch")
		client.MsgCh <- msg
	}
}

func NewClientGroup(num int) *ClientGroup {
	g := &ClientGroup{
		Clients: make([]*Client, num),
	}
	for i := range g.Clients {
		g.Clients[i] = NewClient(-1)
		g.Clients[i].AttachToGroup(g)
	}
	return g
}

func NewClient(userid int64) *Client {
	return &Client{
		Userid: userid,
		MsgCh:  make(chan *string, 10),
	}
}

func (c *Client) AttachToGroup(g *ClientGroup) {
	c.Mgr = g
}

var pool = pkg.NewMsgPool()

func (c *Client) Connect(cometAddr string) {
	// 连接到WebSocket服务
	conn, _, err := websocket.DefaultDialer.Dial(cometAddr, nil)
	if err != nil {
		log.Fatal("连接错误:", err)
		return
	}
	//log.Println("连接到comet:", cometAddr)
	c.conn = pkg.NewConnWs(conn, pool)
	go c.Receive()
	go c.Send()
}

func (c *Client) Auth(str string) {
	// 发送消息到服务器
	token := &msg.CAuthReq{
		Params: []string{str},
		Userid: c.Userid,
	}
	tokenDate, err := proto.Marshal(token)
	if err != nil {
		//log.Println("faild to encode token")
	}
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "auth",
		Seq:  1,
		Data: tokenDate,
	})
	if err != nil {
		//log.Println("发送消息错误:", err)
	}
}

func (c *Client) SendRoomMsg(message *string) {
	sendRoomReq := &msg.CSendRoomReq{
		Userid: c.Userid,
		Roomid: c.Roomids[0],
		Msg:    *message,
	}
	data, err := proto.Marshal(sendRoomReq)
	if err != nil {
		//log.Println("faild to proto marshal", err)
	}
	//log.Println("发送房间消息")
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "sendRoom",
		Seq:  1,
		Data: data,
	})
	if err != nil {
		//log.Println("send err", err)
	}
}

func (c *Client) QueryRoom() {
	cRoomidReq := &msg.CRoomidReq{
		Userid: c.Userid,
	}
	data, err := proto.Marshal(cRoomidReq)
	if err != nil {
		//log.Println("faild to marhal")
		return
	}

	err = c.conn.WriteMsg(&msg.Msg{
		Path: "queryRoom",
		Seq:  1,
		Data: data,
	})
	if err != nil {
		//log.Println("send err:", err)
	}
}

func (c *Client) NewUser() {
	data, err := proto.Marshal(&msg.CNewUserReq{})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "newUser",
		Seq:  1,
		Data: data,
	})
	if err != nil {
		//log.Println("send err:", err)
	}
}

func (c *Client) JoinRoom(roomid int64) {
	data, err := proto.Marshal(&msg.CJoinRoomReq{
		Userid: c.Userid,
		Roomid: roomid,
	})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "joinRoom",
		Seq:  2,
		Data: data,
	})
	if err != nil {
		//log.Println("send err:", err)
	}
	c.Roomids = append(c.Roomids, roomid)
}

func (c *Client) Online() {
	data, err := proto.Marshal(&msg.COnlineReq{
		Userid: c.Userid,
	})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.conn.WriteMsg(&msg.Msg{
		Path: "online",
		Seq:  2,
		Data: data,
	})
	if err != nil {
		//log.Println("send err:", err)
	}
}

func (c *Client) Send() {
	for msg := range c.MsgCh {
		//log.Println("out ch")
		c.SendRoomMsg(msg)
	}
}

func (c *Client) Receive() {
	for {
		message, err := c.conn.ReadMsg()
		if err != nil {
			//log.Println("comet error:", err)
			return
		}
		// //log.Printf("comet reply: %s", message.String())
		switch message.Path {
		case "newUser":
			rt := &msg.CNewUserResp{}
			err := proto.Unmarshal(message.Data, rt)
			if err != nil {
				//log.Println("faild proto", message.Path, err)
				continue
			}
			c.Userid = rt.Userid
			//log.Println("newUser:", c.Userid)
			if c.Mgr != nil {
				c.Mgr.Wait.Done()
			}
		case "auth":
			//log.Println("auth:", c.Userid, string(message.Data))
			c.QueryRoom()
		case "joinRoom":
			//log.Println("joinRoom:", c.Userid, string(message.Data))
			if c.Mgr != nil {
				c.Mgr.Wait.Done()
			}
		case "online":
			//log.Println("online:", c.Userid, string(message.Data))
			if c.Mgr != nil {
				c.Mgr.Wait.Done()
			}
		case "sendRoom":
			//log.Println("send room:", string(message.Data))
		case "queryRoom":
			roomResp := &msg.CRoomidResp{}
			err := proto.Unmarshal(message.Data, roomResp)
			if err != nil {
				//log.Println("faild proto", message.Path, err)
				continue
			}
			c.Roomids = roomResp.Roomid
			//log.Println("query room:", c.Roomids[0])
		case "roomMsg":
			c.ReceiveCount++
			//log.Printf("ch.id[%d] roomMsg receive:%s\n", c.Userid, string(message.Data))
		}
	}
}
