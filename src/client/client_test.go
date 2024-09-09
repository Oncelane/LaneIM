package client_test

import (
	"laneIM/src/client"
	"log"
	"testing"
	"time"
)

func TestManyUser(t *testing.T) {
	num := 1000
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws"}
	// var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)

	{
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
	{
		g.Wait.Add(num)
		for i, c := range g.Clients {
			go func(i int, c *client.Client) {
				c.JoinRoom(1005)
			}(i, c)
		}
		g.Wait.Wait()
		log.Printf("all %d join room", num)
	}
	{
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
	go func() {
		for {

			sum := 0
			for _, c := range g.Clients {
				sum += c.ReceiveCount
			}
			time.Sleep(time.Millisecond)
			if sum == num*num {
				log.Println("receve message count:", sum, "spand time", time.Since(timeStart))
				break
			}
		}
	}()
	select {}
}
