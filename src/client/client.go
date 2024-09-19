package client

import (
	"fmt"
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"log"
	"sync"
	"time"

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
	Room         []int64
	Seq          int64
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
	Num         int
	Clients     []*Client
	Wait        sync.WaitGroup
	MsgCount    int
	TargetCount int
}

func (c *ClientGroup) Send(msg *string) {
	c.TargetCount += c.Num * c.Num
	for _, client := range c.Clients {
		//laneLog.Logger.Infoln("in ch")
		client.MsgCh <- msg
	}
	timeStart := time.Now()
	for {
		sum := 0
		for _, c := range c.Clients {
			sum += c.ReceiveCount
		}
		time.Sleep(time.Millisecond * 1)
		// laneLog.Logger.Debugln("receve count now=", sum, "expect =", c.TargetCount)
		if sum == c.TargetCount {
			laneLog.Logger.Infoln("receve message count: ", sum, " spand time ", time.Since(timeStart))
			break
		}
	}
}

func NewClientGroup(num int) *ClientGroup {
	g := &ClientGroup{
		Clients:  make([]*Client, num),
		MsgCount: 0,
		Num:      num,
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

func (c *Client) SendCometSingle(path string, data []byte) error {
	err := c.conn.WriteMsg(&msg.MsgBatch{
		Msgs: []*msg.Msg{
			{
				Path: path,
				Seq:  c.Seq,
				Data: data,
			},
		},
	})
	c.Seq++
	return err
}

func (c *Client) SendCometBatch(paths []string, datas [][]byte) error {
	if len(paths) != len(datas) {
		return fmt.Errorf("wrong send format, paths.length != datas.lenght")
	}
	msgs := make([]*msg.Msg, len(paths))
	for i := range paths {
		msgs[i].Path = paths[i]
		msgs[i].Data = datas[i]
		msgs[i].Seq = c.Seq
		c.Seq++
	}
	err := c.conn.WriteMsg(&msg.MsgBatch{
		Msgs: msgs,
	})
	return err
}

func (c *Client) Connect(cometAddr string) {
	// 连接到WebSocket服务
	conn, _, err := websocket.DefaultDialer.Dial(cometAddr, nil)
	if err != nil {
		laneLog.Logger.Fatalln("faild to connect comet:", err)
		return
	}
	// laneLog.Logger.Infoln("连接到comet:", cometAddr)
	c.conn = pkg.NewConnWs(conn, pool)
	if c.Mgr != nil {
		c.Mgr.Wait.Done()
	}
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
		laneLog.Logger.Infoln("faild to encode token")
	}
	err = c.SendCometSingle("auth", tokenDate)
	if err != nil {
		laneLog.Logger.Infoln("发送消息错误:", err)
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
		laneLog.Logger.Infoln("faild to proto marshal", err)
	}
	//laneLog.Logger.Infoln("发送房间消息")
	err = c.SendCometSingle("sendRoom", data)
	if err != nil {
		laneLog.Logger.Infoln("send err", err)
	}
}

func (c *Client) QueryRoom() {
	cRoomidReq := &msg.CRoomidReq{
		Userid: c.Userid,
	}
	data, err := proto.Marshal(cRoomidReq)
	if err != nil {
		laneLog.Logger.Infoln("faild to marhal")
		return
	}
	err = c.SendCometSingle("queryRoom", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
	}
}

func (c *Client) NewUser() {
	data, err := proto.Marshal(&msg.CNewUserReq{})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.SendCometSingle("newUser", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
	}
}

func (c *Client) NewRoom() {
	data, err := proto.Marshal(&msg.CNewRoomReq{
		Userid: c.Userid,
	})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.SendCometSingle("newRoom", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
	}
}

func (c *Client) JoinRoom(roomid int64) {
	data, err := proto.Marshal(&msg.CJoinRoomReq{
		Userid: c.Userid,
		Roomid: roomid,
	})
	// laneLog.Logger.Warnf("client[%d] join roomid%d", c.Userid, roomid)
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.SendCometSingle("joinRoom", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
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
	err = c.SendCometSingle("online", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
	}
	// laneLog.Logger.Debugln("send online message success", c.Userid)
}

func (c *Client) Offline() {
	data, err := proto.Marshal(&msg.COfflineReq{
		Userid: c.Userid,
	})
	if err != nil {
		log.Panicln("faild to encode")
	}
	err = c.SendCometSingle("offline", data)
	if err != nil {
		laneLog.Logger.Infoln("send err:", err)
	}
}

func (c *Client) Send() {
	for msg := range c.MsgCh {
		c.SendRoomMsg(msg)
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Receive() {
	for {
		message, err := c.conn.ReadMsg()
		if err != nil {
			if websocket.IsCloseError(err, 1000) {
				// laneLog.Logger.Infoln("comet close normal:", err)
			} else {
				// laneLog.Logger.Infoln("comet close error:", err)
			}
			c.conn.Close()
			return
		}
		for i := range message.Msgs {
			// laneLog.Logger.Infof("comet reply: %s", message.Msgs[i].GetPath())
			switch message.Msgs[i].Path {
			case "newUser":
				rt := &msg.CNewUserResp{}
				err := proto.Unmarshal(message.Msgs[i].Data, rt)
				if err != nil {
					laneLog.Logger.Infoln("faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Userid = rt.Userid
				//laneLog.Logger.Infoln("newUser:", c.Userid)
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
			case "newRoom":
				rt := &msg.CNewRoomResp{}
				err := proto.Unmarshal(message.Msgs[i].Data, rt)
				if err != nil {
					laneLog.Logger.Infoln("faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Room = append(c.Room, rt.Roomid)
				laneLog.Logger.Infoln("client [0] get roomid", rt.Roomid)
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
			case "auth":
				//laneLog.Logger.Infoln("auth:", c.Userid, string(message.Msgs[i].Data))
				c.QueryRoom()
			case "joinRoom":
				//laneLog.Logger.Infoln("joinRoom:", c.Userid, string(message.Msgs[i].Data))
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
			case "online":
				// laneLog.Logger.Infoln("online:", c.Userid, string(message.Msgs[i].Data))
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
			case "sendRoom":
				//laneLog.Logger.Infoln("send room:", string(message.Msgs[i].Data))
			case "queryRoom":
				roomResp := &msg.CRoomidResp{}
				err := proto.Unmarshal(message.Msgs[i].Data, roomResp)
				if err != nil {
					laneLog.Logger.Infoln("faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Roomids = roomResp.Roomid
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
				// laneLog.Logger.Infoln("query room:", c.Roomids[0])
			case "receive":

				c.ReceiveCount++
				// laneLog.Logger.Infof("ch.id[%d] roomMsg receive:%s\n", c.Userid, string(message.Msgs[i].Data))
			case "offline":
				if c.Mgr != nil {
					c.Mgr.Wait.Done()
				}
			}
		}
	}
}
