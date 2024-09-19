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

var num int = 1000

func TestManyUser(t *testing.T) {
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)

	{ // connetc
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				// block
				c.Connect(cometAddr[i%len(cometAddr)])

			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d connect and newuser time %v", num, time.Since(start))
	}
	{ // new user
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.NewUser()
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d connect and newuser time %v", num, time.Since(start))
	}
	{ // new room
		start := time.Now()
		g.Wait.Add(1)
		go func() {
			g.Clients[0].NewRoom()
		}()
		g.Wait.Wait()
		laneLog.Logger.Infoln("new roomid:", g.Clients[0].Room[0], " spand time ", time.Since(start))
	}

	{ // join room
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.JoinRoom(g.Clients[0].Room[0])
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d join room spand time %v", num, time.Since(start))
	}

	// save user to disk
	userids := make([]int64, len(g.Clients))
	{
		for i, c := range g.Clients {
			userids[i] = c.Userid
		}
		file, err := os.Create("userids")
		if err != nil {
			laneLog.Logger.Fatalln("save error", err)
			t.Error(err)
		}
		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)
		e.Encode(userids)
		file.Write(b.Bytes())
		laneLog.Logger.Infof("all %d user id save", num)
	}

	{ // set online
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.Online()
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d online spand time %v", num, time.Since(start))
	}

	msg := "hello"
	g.Send(&msg)
	msg = "22222"
	g.Send(&msg)
	msg = "33333"
	g.Send(&msg)
	msg = "我可不觉得这段话很长，算是一般长度"
	g.Send(&msg)
	{ // set offline
		start := time.Now()
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.Offline()
			}(i, c)
		}
		laneLog.Logger.Infof("all %d offline spand time %v", num, time.Since(start))
	}
	time.Sleep(time.Second * 2)
	laneLog.Logger.Infoln("end")
}
func TestCacheUser(t *testing.T) {
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)
	{ // read user from disk
		userids := make([]int64, len(g.Clients))
		data, err := os.ReadFile("userids")
		if err != nil {
			laneLog.Logger.Fatalln("save error", err)
			t.Error(err)
		}
		b := bytes.NewBuffer(data)
		e := gob.NewDecoder(b)
		e.Decode(&userids)
		for i, c := range g.Clients {
			c.Userid = userids[i]
			// laneLog.Logger.Infoln("read userids ", userids[i])
		}
		laneLog.Logger.Infoln("all read userids success")
	}

	{ // connetc user
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				// block
				c.Connect(cometAddr[i%len(cometAddr)])
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d connect spand time %v", num, time.Since(start))
	}

	{ // query room
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.QueryRoom()
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d query room spand time %v", num, time.Since(start))
	}

	{ // set online
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.Online()
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d online spand time %v", num, time.Since(start))
	}

	msg := "hello"
	g.Send(&msg)
	msg = "22222"
	g.Send(&msg)
	msg = "33333"
	g.Send(&msg)
	{ // set offline
		start := time.Now()
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.Offline()
			}(i, c)
		}
		laneLog.Logger.Infof("all %d offline spand time %v", num, time.Since(start))
	}
	time.Sleep(time.Second * 2)
	laneLog.Logger.Infoln("end")
}
