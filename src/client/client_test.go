package client_test

import (
	"bytes"
	"encoding/gob"
	"laneIM/src/client"
	"laneIM/src/pkg/laneLog"
	"os"
	"strconv"
	"testing"
	"time"
)

var cometAddr []string = []string{"ws://127.0.0.1:40050/ws"}

// var cometAddr []string = []string{"ws://172.29.178.158:40050/ws"}

// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
var snum = 2

func TestSimulate(t *testing.T) {
	g := client.NewClientGroup(snum)

	g.Wait.Add(snum)
	for _, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.Connect(cometAddr[0])
			// laneLog.Logger.Debugln("pass1")
			c.NewUser()
			// laneLog.Logger.Debugln("pass2")
			c.Online()
			// laneLog.Logger.Debugln("pass3")
		}()
	}
	g.Wait.Wait()
	// select {}
	g.Wait.Add(snum)
	roomid := g.Clients[0].NewRoom()
	for _, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.JoinRoom(roomid)
		}()
	}
	g.Wait.Wait()

	for _, c := range g.Clients {
		go func() {
			c.SendRoomMsg(roomid, "1")
		}()
	}

	time.Sleep(time.Second * 2)
	receiveBytes := g.ReceiveBytes()
	laneLog.Logger.Infoln("recevie bytes:", receiveBytes)
	g.Wait.Add(snum / 2)
	for i := 0; i < snum/2; i++ {
		go func() {
			defer g.Wait.Done()

			g.Clients[i].Subon(roomid)

		}()
	}
	g.Wait.Wait()

	for _, c := range g.Clients {
		go func() {
			c.SendRoomMsg(roomid, "2")
		}()
	}
	time.Sleep(time.Second * 2)
	receiveBytes = g.ReceiveBytes() - receiveBytes
	laneLog.Logger.Infoln("recevie bytes:", receiveBytes)
}

var tnum = 5000

// 创建10000个用户，并加入新房间

func TestCreateUsersAndJoinNewRoom(t *testing.T) {
	g := client.NewClientGroup(tnum)

	var limit = 1000
	i := 0
	userids := make([]int64, tnum)
	for {
		left := tnum - i
		var end int
		if left > limit {
			end = i + limit
		} else if left == 0 {
			break
		} else {
			end = i + left
		}
		var count = end - i
		g.Wait.Add(count)

		for ; i < end; i++ {
			go func(i int) {
				defer g.Wait.Done()
				g.Clients[i].Connect(cometAddr[0])
				userids[i] = g.Clients[i].NewUser()
				// c.Online()
			}(i)
		}
		g.Wait.Wait()
		laneLog.Logger.Infoln("new user", count)
	}

	roomid := g.Clients[0].NewRoom()
	laneLog.Logger.Infoln("roomid = ", roomid)
	i = 0
	for {
		left := tnum - i
		var end int
		if left > limit {
			end = i + limit
		} else if left == 0 {
			break
		} else {
			end = i + left
		}
		var count = end - i
		g.Wait.Add(count)

		for ; i < end; i++ {
			go func(i int) {
				defer g.Wait.Done()
				g.Clients[i].JoinRoom(roomid)
				// c.Online()
			}(i)
		}
		g.Wait.Wait()
		laneLog.Logger.Infoln("join room", count)
	}

	{ // save user to disk
		file, err := os.Create("userids")
		if err != nil {
			laneLog.Logger.Fatalln("save error", err)
			t.Error(err)
		}
		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)
		e.Encode(roomid)
		e.Encode(userids)
		file.Write(b.Bytes())
		laneLog.Logger.Infoln("roomid and userids save in disk")
	}
}

// 创建10000个用户，并加入指定房间
func TestCreateUsersAndJoinGivenRoom(t *testing.T) {

	var roomid int64 = 1837561315430760448

	g := client.NewClientGroup(tnum)

	var limit = 1000
	i := 0
	userids := make([]int64, tnum)
	for {
		left := tnum - i
		var end int
		if left > limit {
			end = i + limit
		} else if left == 0 {
			break
		} else {
			end = i + left
		}
		var count = end - i
		g.Wait.Add(count)

		for ; i < end; i++ {
			go func(i int) {
				defer g.Wait.Done()
				g.Clients[i].Connect(cometAddr[0])
				userids[i] = g.Clients[i].NewUser()
				// c.Online()
			}(i)
		}
		g.Wait.Wait()
		laneLog.Logger.Infoln("new user", count)
	}

	i = 0
	for {
		left := tnum - i
		var end int
		if left > limit {
			end = i + limit
		} else if left == 0 {
			break
		} else {
			end = i + left
		}
		var count = end - i
		g.Wait.Add(count)

		for ; i < end; i++ {
			go func(i int) {
				defer g.Wait.Done()
				g.Clients[i].JoinRoom(roomid)
				// c.Online()
			}(i)
		}
		g.Wait.Wait()
		laneLog.Logger.Infoln("join room", count)
	}

	{ // save user to disk
		file, err := os.Create("userids")
		if err != nil {
			laneLog.Logger.Fatalln("save error", err)
			t.Error(err)
		}
		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)
		e.Encode(roomid)
		e.Encode(userids)
		file.Write(b.Bytes())
		laneLog.Logger.Infoln("roomid and userids save in disk")
	}
}

func GetUseridAndRoomidFromDisk() (int64, []int64) {
	// read user from disk
	userids := make([]int64, tnum)
	var roomid int64
	data, err := os.ReadFile("userids")
	if err != nil {
		laneLog.Logger.Fatalln("save error", err)
	}
	b := bytes.NewBuffer(data)
	e := gob.NewDecoder(b)
	e.Decode(&roomid)
	e.Decode(&userids)
	laneLog.Logger.Infoln("read userids success")
	return roomid, userids
}

func TestOneRoomConstentlySend(t *testing.T) {
	g := client.NewClientGroup(tnum)

	roomid, userids := GetUseridAndRoomidFromDisk()
	laneLog.Logger.Infoln("roomid:", roomid)
	for i, c := range g.Clients {
		c.Userid = userids[i]
	}
	laneLog.Logger.Infoln("init userids")

	g.Wait.Add(tnum)
	for _, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.Connect(cometAddr[0])
			c.Online()
			c.Subon(roomid)
		}()
	}
	g.Wait.Wait()
	laneLog.Logger.Infoln("Subon userids")

	laneLog.Logger.Infof("等待%d秒后开始传输", (60 - time.Now().Second()))
	time.Sleep(time.Second * time.Duration((60 - time.Now().Second())))

	start := time.Now()
	done := make(chan struct{})
	sendCount := 0
	go func() {
		for {
			if time.Since(start).Seconds() > 30 {
				done <- struct{}{}
				break
			}
			time.Sleep(time.Millisecond * 100)

			sendUser := 300
			g.Wait.Add(sendUser)
			// laneLog.Logger.Infoln("send count =", sendUser)
			for i := range sendUser {
				go func() {
					defer g.Wait.Done()
					sendCount++
					g.Clients[i].SendRoomMsg(roomid, "hello")
				}()
			}
			g.Wait.Wait()
			// laneLog.Logger.Infoln("send ack =", sendUser, "averavg ack =", g.AverageAck(sendUser))
		}
	}()
	go g.ReceiveCount(time.Second)
	go func() {
		lastCount := 0
		for {
			time.Sleep(time.Second)
			laneLog.Logger.Infoln("send count=", sendCount-lastCount)
			lastCount = sendCount
		}
	}()
	<-done
	for _, c := range g.Clients {
		go func() {
			c.Offline()
		}()
	}
	time.Sleep(time.Second * 2)
}

func TestSimulateMany(t *testing.T) {
	g := client.NewClientGroup(tnum)

	g.Wait.Add(tnum)
	for i, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.Connect(cometAddr[i%2])
			// c.Online()
		}()
	}
	g.Wait.Wait()
	roomid, userids := GetUseridAndRoomidFromDisk()
	for i, c := range g.Clients {
		c.Userid = userids[i]
	}
	laneLog.Logger.Infoln("init userids")

	for _, c := range g.Clients {
		time.Sleep(time.Millisecond * 20)
		c.SendRoomMsg(roomid, "hello")
	}
	laneLog.Logger.Infoln("subscibe roomid")
	for _, c := range g.Clients {
		time.Sleep(time.Millisecond * 20)
		c.SendRoomMsg(roomid, "hello")
	}

	time.Sleep(time.Second * 2)
	sendBytes := g.SendBytes()
	laneLog.Logger.Infoln("send bytes:", sendBytes)

	receiveBytes := g.ReceiveBytes()
	laneLog.Logger.Infoln("recevie bytes:", receiveBytes)

	laneLog.Logger.Infof("num of %d user subscribe room %d", tnum/2, roomid)
	g.Wait.Add(tnum / 2)
	for i := 0; i < tnum/2; i++ {
		go func() {
			defer g.Wait.Done()

			g.Clients[i].Subon(roomid)

		}()
	}
	g.Wait.Wait()

	for _, c := range g.Clients {
		go func() {
			c.SendRoomMsg(roomid, "hello")
		}()
	}
	time.Sleep(time.Second * 2)

	sendBytes = g.SendBytes() - sendBytes
	laneLog.Logger.Infoln("send bytes:", sendBytes)

	receiveBytes = g.ReceiveBytes() - receiveBytes
	laneLog.Logger.Infoln("recevie bytes:", receiveBytes)

}

func InitOneClientGroup(b *testing.B, bnum int) (*client.ClientGroup, int64) {
	b.Helper()
	g := client.NewClientGroup(bnum)

	g.Wait.Add(bnum)
	for _, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.Connect(cometAddr[0])
			c.NewUser()
			c.Online()
		}()
	}
	g.Wait.Wait()
	// select {}
	g.Wait.Add(bnum)
	roomid := g.Clients[0].NewRoom()
	for _, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.JoinRoom(roomid)
		}()
	}
	g.Wait.Wait()

	g.Wait.Add(bnum)
	for i := 0; i < bnum; i++ {
		go func() {
			defer g.Wait.Done()
			g.Clients[i].Subon(roomid)

		}()
	}
	g.Wait.Wait()
	return g, roomid
}

func BenchmarkOneClientSend(b *testing.B) {
	var bnum = 1000
	g, roomid := InitOneClientGroup(b, bnum)
	start := time.Now()
	for range b.N {
		g.Clients[0].SendRoomMsg(roomid, "hello")
	}
	g.WaitMessageCount(bnum * b.N)
	laneLog.Logger.Infof("send %d,spand time %v", b.N, time.Since(start))
}

// generator 100000 message
func HelpGeneratorMessage(b *testing.B, groupSize int, messageSize int) (*client.ClientGroup, int64) {
	b.Helper()
	g, roomid := InitOneClientGroup(b, groupSize)
	for i := range messageSize {
		g.Clients[0].SendRoomMsg(roomid, "this is message "+strconv.Itoa(i))
	}
	g.WaitMessageCount(groupSize * messageSize)
	return g, roomid
}

func BenchmarkPaging(b *testing.B) {
	messageSize := 1000
	roomSize := 100
	pageSize := 20
	g, roomid := HelpGeneratorMessage(b, roomSize, messageSize)
	laneLog.Logger.Infoln("save ", messageSize, " messages")
	last, err := g.Clients[0].QueryLast(roomid)
	g.Clients[0].LastMessageId = last.MessageId
	g.Clients[0].LastMessageUnix = last.TimeUnix.AsTime().Local()
	if err != nil {
		b.Fatal(err)
	}
	laneLog.Logger.Infoln("query last =", last.String())
	var pageIndex = 0
	c := g.Clients[0]
	for {

		msgs, err := c.QueryPaging(roomid, c.LastMessageId, c.LastMessageUnix, int64(pageSize))
		if err != nil {
			b.Fatal(err)
		}
		if len(msgs.Msgs) < pageSize {
			laneLog.Logger.Infoln("===========end%d==========", pageIndex)
			break
		}
		lastMsg := msgs.Msgs[len(msgs.Msgs)-1]
		laneLog.Logger.Infof("=========page%d=>lastMsg:%s", pageIndex, string(lastMsg.Data))
		c.SetQueryLastIndex(lastMsg.Messageid, lastMsg.Timeunix.AsTime().Local())
		pageIndex++
	}

}

// 测试上线后获取离线消息
// func TestPageChatMessageUnReadCountByMessageid(t *testing.T) {
// 	db :=
// 	db.PageChatMessageUnReadCountByMessageid()
// }
