package client

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	//control
	Mgr           *ClientGroup
	conn          *pkg.ConnWs
	Userid        int64
	GroupSeq      int
	Seq           int64
	Roomids       []int64
	LastMessageId int64

	// statistics
	SendByte         int64
	ReceiveByte      int64
	RoomReceiveBytes map[int64]int64

	ReceiveCount int
	ackSpandTime time.Duration

	// sync
	waitNewRoom       chan int64
	waitAuth          chan string
	waitSendAck       chan string
	waitQueryRoom     chan []int64
	waitNewUser       chan int64
	waitJoinRoom      chan string
	waitOnline        chan string
	waitOffline       chan string
	waitLastMessageid chan int64
	waitPaging        chan *msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs
	waitSubOn         chan string
	waitSubOff        chan string
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

func (c *ClientGroup) WaitMessageCount(count int) {
	c.TargetCount += count
	timeStart := time.Now()
	for {
		sum := 0
		for _, c := range c.Clients {
			sum += c.ReceiveCount
		}
		time.Sleep(time.Millisecond)
		if sum == c.TargetCount {
			laneLog.Logger.Infoln("[client] receve message sum: ", sum, "target", c.TargetCount, " spand time ", time.Since(timeStart))
			break
		}
	}
}

func (c *ClientGroup) ReceiveCount(interval time.Duration) {
	var lastBytes int64 = 0
	sum := 0
	for {
		var tmp = 0
		for _, c := range c.Clients {
			tmp += c.ReceiveCount
		}

		receiveByte := c.ReceiveBytes()
		laneLog.Logger.Infoln("[client] receve message count: ", tmp-sum, "bytes:", receiveByte-lastBytes)
		sum = tmp
		lastBytes = receiveByte
		time.Sleep(interval)
	}
}

func (c *ClientGroup) ReceiveBytes() (totalBytes int64) {
	for _, cl := range c.Clients {
		totalBytes += cl.ReceiveBytes()
	}
	return totalBytes
}

func (c *ClientGroup) AverageAck(num int) (averageSpandMillisecond int64) {
	var averageSpand time.Duration
	for i := 0; i < num; i++ {
		averageSpand += c.Clients[i].ackSpandTime
	}
	return averageSpand.Milliseconds() / int64(num)
}

func (c *Client) ReceiveBytes() int64 {
	return c.ReceiveByte
}
func (c *ClientGroup) SendBytes() (totalBytes int64) {
	for _, cl := range c.Clients {
		totalBytes += cl.SendBytes()
	}
	return totalBytes
}

func (c *Client) SendBytes() int64 {
	return c.SendByte
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
		g.Clients[i].GroupSeq = i
	}
	return g
}

func NewClient(userid int64) *Client {
	return &Client{
		Userid:            userid,
		waitNewRoom:       make(chan int64),
		waitAuth:          make(chan string),
		waitSendAck:       make(chan string),
		waitQueryRoom:     make(chan []int64),
		waitNewUser:       make(chan int64),
		waitJoinRoom:      make(chan string),
		waitOnline:        make(chan string),
		waitOffline:       make(chan string),
		waitLastMessageid: make(chan int64),
		waitPaging:        make(chan *msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs),
		waitSubOn:         make(chan string),
		waitSubOff:        make(chan string),
		RoomReceiveBytes:  make(map[int64]int64),
	}
}

func (c *Client) AttachToGroup(g *ClientGroup) {
	c.Mgr = g
}

var pool = pkg.NewMsgPool()

func (c *Client) sendCometSingle(path string, data []byte) error {
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

func (c *Client) SendCometBatch(in []*msg.Msg) error {
	err := c.conn.WriteMsg(&msg.MsgBatch{
		Msgs: in,
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
	go c.Receive()
}

func (c *Client) Auth(str string) string {
	// 发送消息到服务器
	token := &msg.CAuthReq{
		Params: []string{str},
		Userid: c.Userid,
	}
	tokenDate, err := proto.Marshal(token)
	if err != nil {
		laneLog.Logger.Infoln("[client] faild to encode token")
	}
	err = c.sendCometSingle("auth", tokenDate)
	if err != nil {
		laneLog.Logger.Fatalln("[client] 发送消息错误:", err)
	}
	return <-c.waitAuth
}

func (c *Client) SendRoomMsg(roomid int64, message string) string {

	sendRoomReq := &msg.CSendRoomReq{
		Userid: c.Userid,
		Roomid: roomid,
		Msg:    message,
	}
	data, err := proto.Marshal(sendRoomReq)
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to proto marshal", err)
	}
	//laneLog.Logger.Infoln("[client] 发送房间消息")
	sendTime := time.Now()

	err = c.sendCometSingle("sendRoom", data)
	if err != nil {
		laneLog.Logger.Infoln("[client] send err", err)
	}
	c.SendByte += int64(len(data))
	s := <-c.waitSendAck
	c.ackSpandTime = time.Since(sendTime)
	return s
}
func (c *Client) Subon(roomid int64) string {

	err := c.sendCometSingle("subon", []byte(strconv.FormatInt(roomid, 10)))
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	return <-c.waitSubOn
}
func (c *Client) Suboff(roomid int64) string {

	err := c.sendCometSingle("suboff", []byte(strconv.FormatInt(roomid, 10)))
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	return <-c.waitSubOff
}

func (c *Client) QueryRoom() []int64 {
	cRoomidReq := &msg.CRoomidReq{
		Userid: c.Userid,
	}
	data, err := proto.Marshal(cRoomidReq)
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to marhal")
		return nil
	}
	err = c.sendCometSingle("queryRoom", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	return <-c.waitQueryRoom
}

func (c *Client) NewUser() int64 {
	data, err := proto.Marshal(&msg.CNewUserReq{})
	if err != nil {
		laneLog.Logger.Fatalln("faild to encode")
	}
	err = c.sendCometSingle("newUser", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	return <-c.waitNewUser
}

func (c *Client) LoadUser(userid int64) {
	c.Userid = userid
}

func (c *Client) NewRoom() int64 {
	data, err := proto.Marshal(&msg.CNewRoomReq{
		Userid: c.Userid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to encode")
	}
	err = c.sendCometSingle("newRoom", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	return <-c.waitNewRoom
}

func (c *Client) JoinRoom(roomid int64) string {
	data, err := proto.Marshal(&msg.CJoinRoomReq{
		Userid: c.Userid,
		Roomid: roomid,
	})
	// laneLog.Logger.Warnf("client[%d] join roomid%d", c.Userid, roomid)
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to encode")
	}
	err = c.sendCometSingle("joinRoom", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	c.Roomids = append(c.Roomids, roomid)
	return <-c.waitJoinRoom
}

func (c *Client) Online() string {
	data, err := proto.Marshal(&msg.COnlineReq{
		Userid: c.Userid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to encode")
	}
	err = c.sendCometSingle("online", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	// laneLog.Logger.Debugln("send online message success", c.Userid)
	return <-c.waitOnline
}

func (c *Client) Offline() {
	data, err := proto.Marshal(&msg.COfflineReq{
		Userid: c.Userid,
	})
	if err != nil {
		laneLog.Logger.Fatalln("[client] faild to encode")
	}
	err = c.sendCometSingle("offline", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
	}
	// 被动关闭，协商让服务器关闭
	// c.conn.ForceClose()
	// return <-c.waitOffline
}

func (c *Client) QueryLastMessageId(roomid int64) (int64, error) {
	err := c.sendCometSingle("last", []byte(strconv.FormatInt(roomid, 10)))
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
		return 0, err
	}
	return <-c.waitLastMessageid, nil
}
func (c *Client) QueryPaging(roomid int64, lastMsgId int64, limit int64) (*msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs, error) {
	data, err := proto.Marshal(&msg.CQueryStoreMessageReq{
		Roomid:    roomid,
		MessageId: lastMsgId,
		Size:      limit,
	})
	if err != nil {
		laneLog.Logger.Fatalln(err)
	}
	err = c.sendCometSingle("his", data)
	if err != nil {
		laneLog.Logger.Fatalln("[client] send err:", err)
		return nil, err
	}
	return <-c.waitPaging, nil
}

func (c *Client) PassiveClose() error {
	return c.conn.ForceClose()
}

func (c *Client) ForceClose() error {
	return c.conn.ForceClose()
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
			c.ForceClose()
			return
		}
		for i, m := range message.Msgs {
			// laneLog.Logger.Infof("comet reply: %s", message.Msgs[i].GetPath())
			switch m.Path {
			case "newUser":
				rt := &msg.CNewUserResp{}
				err := proto.Unmarshal(m.Data, rt)
				if err != nil {
					laneLog.Logger.Fatalln("[client] faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Userid = rt.Userid
				c.waitNewUser <- rt.Userid
				//laneLog.Logger.Infoln("newUser:", c.Userid)
			case "newRoom":
				rt := &msg.CNewRoomResp{}
				err := proto.Unmarshal(message.Msgs[i].Data, rt)
				if err != nil {
					laneLog.Logger.Fatalln("[client] faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Roomids = append(c.Roomids, rt.Roomid)
				laneLog.Logger.Infoln("[client] get roomid", rt.Roomid)
				c.waitNewRoom <- rt.Roomid
			case "auth":
				//laneLog.Logger.Infoln("auth:", c.Userid, string(message.Msgs[i].Data))
				// c.QueryRoom()
				c.waitAuth <- string(m.Data)
			case "joinRoom":
				//laneLog.Logger.Infoln("joinRoom:", c.Userid, string(message.Msgs[i].Data))
				c.waitJoinRoom <- string(m.Data)
			case "online":
				// laneLog.Logger.Infoln("online:", c.Userid, string(message.Msgs[i].Data))
				c.waitOnline <- string(m.Data)
			case "sendRoom":
				//laneLog.Logger.Infoln("send room:", string(message.Msgs[i].Data))

				c.waitSendAck <- string(m.Data)
			case "queryRoom":
				roomResp := &msg.CRoomidResp{}
				err := proto.Unmarshal(message.Msgs[i].Data, roomResp)
				if err != nil {
					laneLog.Logger.Fatalln("[client] faild proto", message.Msgs[i].Path, err)
					continue
				}
				c.Roomids = roomResp.Roomid
				c.waitQueryRoom <- roomResp.Roomid
				laneLog.Logger.Infoln("query room:", c.Roomids[0])

			case "offline":
				// c.waitOffline <- string(m.Data)
			case "last":
				msgId, err := strconv.ParseInt(string(message.Msgs[i].Data), 10, 64)
				if err != nil {
					laneLog.Logger.Fatalln(err)
				}
				c.LastMessageId = msgId
				c.waitLastMessageid <- msgId
			case "his":
				msgs := new(msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs)
				err := proto.Unmarshal(message.Msgs[i].Data, msgs)
				if err != nil {
					laneLog.Logger.Fatalln(err)
				}
				c.waitPaging <- msgs
				// laneLog.Logger.Infof("user[%d] receive his msg {==-%v-==}", c.Userid, msgs.String())
			case "receive": //online message
				// laneLog.Logger.Infof("ch.id[%d] roomMsg receive:%s\n", c.GroupSeq, string(message.Msgs[i].Data))
				c.ReceiveByte += int64(len(m.Data))
				c.ReceiveCount++
			case "onlyCount":
				countMsg := new(msg.COnlyCountMessage)
				err := proto.Unmarshal(m.Data, countMsg)
				if err != nil {
					laneLog.Logger.Fatalln(err)
				}
				// laneLog.Logger.Infof("[onlyCount] user[%d] roomid[%d], receive count[%d]", c.GroupSeq, countMsg.Roomid, countMsg.Count)
				c.RoomReceiveBytes[countMsg.Roomid] += int64(len(m.Data))
				c.ReceiveByte += int64(len(m.Data))
			case "subon":
				c.waitSubOn <- string(m.Data)
			case "suboff":
				c.waitSubOff <- string(m.Data)
			}
		}
	}
}
