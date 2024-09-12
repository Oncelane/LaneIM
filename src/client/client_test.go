package client_test

import (
	"laneIM/src/client"
	"laneIM/src/pkg/laneLog.go"
	"testing"
	"time"
)

func TestManyUser(t *testing.T) {
	num := 4000
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)

	{ // new user
		start := time.Now()
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				// block
				c.Connect(cometAddr[i%len(cometAddr)])
				// block
				c.NewUser()
			}(i, c)
		}
		g.Wait.Wait()
		laneLog.Logger.Infof("all %d connetc and auth spand time %v", num, time.Since(start))
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
	timeStart := time.Now()

	ch := make(chan struct{})
	go func() {
		for {

			sum := 0
			for _, c := range g.Clients {
				sum += c.ReceiveCount
			}
			time.Sleep(time.Millisecond * 10)
			if sum == num*num {
				laneLog.Logger.Infoln("receve message count: ", sum, " spand time ", time.Since(timeStart))
				ch <- struct{}{}
				break
			}
		}
	}()
	<-ch
	// time.Sleep(time.Second * 30)
	laneLog.Logger.Infoln("start to send second message :hihihi")
	msg = "hihihi"
	g.Send(&msg)
	timeStart = time.Now()
	go func() {
		for {

			sum := 0
			for _, c := range g.Clients {
				sum += c.ReceiveCount
			}
			time.Sleep(time.Millisecond * 10)
			if sum == num*num*2 {
				laneLog.Logger.Infoln("receve message count: ", sum, " spand time ", time.Since(timeStart))
				ch <- struct{}{}
				break
			}
		}
	}()
	<-ch
	laneLog.Logger.Infoln("end")
}
