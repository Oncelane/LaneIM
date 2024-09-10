package client_test

import (
	"laneIM/src/client"
	"log"
	"testing"
	"time"
)

func TestManyUser(t *testing.T) {
	num := 10
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws"}
	// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)

	{ // new user
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
		log.Printf("all %d connetc and auth", num)
	}

	{ // new room
		g.Wait.Add(1)
		go func() {
			g.Clients[0].NewRoom()
		}()
		g.Wait.Wait()
		log.Println("new roomid:", g.Clients[0].Room[0])
	}

	{ // join room
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.JoinRoom(g.Clients[0].Room[0])
			}(i, c)
		}
		g.Wait.Wait()
		log.Printf("all %d join room", num)
	}

	{ // set online
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.Online()
			}(i, c)
		}
		g.Wait.Wait()
		log.Printf("all %d online", num)
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
			time.Sleep(time.Millisecond)
			if sum == num*num {
				log.Println("receve message count:", sum, "spand time", time.Since(timeStart))
				ch <- struct{}{}
				break
			}
		}
	}()
	<-ch
	log.Println("end")
}
