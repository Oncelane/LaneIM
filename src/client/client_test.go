package client_test

import (
	"bytes"
	"encoding/gob"
	"laneIM/src/client"
	"laneIM/src/pkg/laneLog"
	"os"
	"testing"
	"time"
)

// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws"}
var cometAddr []string = []string{"ws://172.29.178.158:40050/ws"}

// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
var snum = 10

func TestSimulate(t *testing.T) {
	g := client.NewClientGroup(snum)

	g.Wait.Add(snum)
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
			c.SendRoomMsg(roomid, "hello")
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
			c.SendRoomMsg(roomid, "hello")
		}()
	}
	time.Sleep(time.Second * 2)
	receiveBytes = g.ReceiveBytes() - receiveBytes
	laneLog.Logger.Infoln("recevie bytes:", receiveBytes)
}

var tnum = 10000

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
	for i, c := range g.Clients {
		go func() {
			defer g.Wait.Done()
			c.Connect(cometAddr[i%2])
			c.Online()
			c.Subon(roomid)
		}()
	}
	g.Wait.Wait()
	laneLog.Logger.Infoln("Subon userids")
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

// func TestManyUser(t *testing.T) {

// 	// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
// 	g := client.NewClientGroup(num)

// 	{ // connetc
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				// block
// 				c.Connect(cometAddr[i%len(cometAddr)])

// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d connect and newuser time %v", num, time.Since(start))
// 	}
// 	{ // new user
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.NewUser()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d connect and newuser time %v", num, time.Since(start))
// 	}
// 	{ // new room
// 		start := time.Now()
// 		g.Wait.Add(1)
// 		go func() {
// 			g.Clients[0].NewRoom()
// 		}()
// 		g.Wait.Wait()
// 		laneLog.Logger.Infoln("[client] new roomid:", g.Clients[0].Roomids[0], " spand time ", time.Since(start))
// 	}

// 	{ // join room
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.JoinRoom(g.Clients[0].Roomids[0])
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d join room spand time %v", num, time.Since(start))
// 	}

// 	// save user to disk
// 	userids := make([]int64, len(g.Clients))
// 	{
// 		for i, c := range g.Clients {
// 			userids[i] = c.Userid
// 		}
// 		file, err := os.Create("userids")
// 		if err != nil {
// 			laneLog.Logger.Fatalln("save error", err)
// 			t.Error(err)
// 		}
// 		b := new(bytes.Buffer)
// 		e := gob.NewEncoder(b)
// 		e.Encode(userids)
// 		file.Write(b.Bytes())
// 		laneLog.Logger.Infof("[client] all %d user id save", num)
// 	}

// 	{ // set online
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Online()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d online spand time %v", num, time.Since(start))
// 	}

// 	// msg := "hello"
// 	// g.Send(&msg)
// 	// msg = "22222"
// 	// g.Send(&msg)
// 	// msg = "33333"
// 	// g.Send(&msg)
// 	// msg = "我可不觉得这段话很长，算是一般长度"
// 	// g.Send(&msg)
// 	{ // set offline
// 		start := time.Now()
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Offline()
// 			}(i, c)
// 		}
// 		laneLog.Logger.Infof("[client] all %d offline spand time %v", num, time.Since(start))
// 	}
// 	time.Sleep(time.Second * 2)
// 	laneLog.Logger.Infoln("[client] end")
// }
// func TestCacheUser(t *testing.T) {
// 	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
// 	g := client.NewClientGroup(num)
// 	{ // read user from disk
// 		userids := make([]int64, len(g.Clients))
// 		data, err := os.ReadFile("userids")
// 		if err != nil {
// 			laneLog.Logger.Fatalln("save error", err)
// 			t.Error(err)
// 		}
// 		b := bytes.NewBuffer(data)
// 		e := gob.NewDecoder(b)
// 		e.Decode(&userids)
// 		for i, c := range g.Clients {
// 			c.Userid = userids[i]
// 			// laneLog.Logger.Infoln("[client] read userids ", userids[i])
// 		}
// 		laneLog.Logger.Infoln("[client] all read userids success")
// 	}

// 	{ // connetc user
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				// block
// 				c.Connect(cometAddr[i%len(cometAddr)])
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d connect spand time %v", num, time.Since(start))
// 	}

// 	{ // query room
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.QueryRoom()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d query room spand time %v", num, time.Since(start))
// 	}

// 	{ // set online
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Online()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d online spand time %v", num, time.Since(start))
// 	}

// 	// msg := "测试长消息的发送延迟，总共一亿条呢，如果换成图片什么的传输的payload将会更大，如果是一个3mb的图片，转发性能如何呢"
// 	// msg := "测试长消息的发送延迟，总共一亿条呢"
// 	// g.Send(&msg)
// 	// msg = "22222"
// 	// g.Send(&msg)
// 	// msg = "testLastMessage"
// 	// g.Send(&msg)
// 	{ // set offline
// 		start := time.Now()
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Offline()
// 			}(i, c)
// 		}
// 		laneLog.Logger.Infof("[client] all %d offline spand time %v", num, time.Since(start))
// 	}
// 	time.Sleep(time.Second * 2)
// 	laneLog.Logger.Infoln("[client] end")
// }

// func TestPageging(t *testing.T) {
// 	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
// 	g := client.NewClientGroup(num)
// 	{ // read user from disk
// 		userids := make([]int64, len(g.Clients))
// 		data, err := os.ReadFile("userids")
// 		if err != nil {
// 			laneLog.Logger.Fatalln("save error", err)
// 			t.Error(err)
// 		}
// 		b := bytes.NewBuffer(data)
// 		e := gob.NewDecoder(b)
// 		e.Decode(&userids)
// 		for i, c := range g.Clients {
// 			c.Userid = userids[i]
// 			// laneLog.Logger.Infoln("read userids ", userids[i])
// 		}
// 		laneLog.Logger.Infoln("[client] all read userids success")
// 	}

// 	{ // connetc user
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				// block
// 				c.Connect(cometAddr[i%len(cometAddr)])
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d connect spand time %v", num, time.Since(start))
// 	}

// 	{ // query room
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.QueryRoom()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d query room spand time %v", num, time.Since(start))
// 	}

// 	{ // set online
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Online()
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d online spand time %v", num, time.Since(start))
// 	}

// 	{ //query lastMessageid
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.QueryLastMessageId(c.Roomids[0])
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d query lastMessageid spand time %v", num, time.Since(start))
// 	}

// 	{ //query lastRoomMessagePage 1
// 		start := time.Now()
// 		g.Wait.Add(num)
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.QueryPaging(c.Roomids[0], c.LastMessageId+1, 100)
// 			}(i, c)
// 		}
// 		g.Wait.Wait()
// 		laneLog.Logger.Infof("[client] all %d query lastRoomMessagePage spand time %v", num, time.Since(start))
// 	}

// 	{ // set offline
// 		start := time.Now()
// 		for i, c := range g.Clients {
// 			go func(i int, c *client.Client) {
// 				c.Offline()
// 			}(i, c)
// 		}
// 		laneLog.Logger.Infof("[client] all %d offline spand time %v", num, time.Since(start))
// 	}
// 	time.Sleep(time.Second * 2)
// 	laneLog.Logger.Infoln("[client] end")
// }

// // 写入
// func TestByte(t *testing.T) {
// 	str := "测试长消息的发送延迟，总共一亿条呢"
// 	laneLog.Logger.Infoln(len(str) * 10000)
// }
