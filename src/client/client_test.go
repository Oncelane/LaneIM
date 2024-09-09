package client_test

import (
	"laneIM/src/client"
	"log"
	"sync"
	"testing"
	"time"
)

func TestManyUser(t *testing.T) {
	num := 2
	var cometAddr []string = []string{"ws://127.0.0.1:40050/ws", "ws://127.0.0.1:40051/ws"}
	g := client.NewClientGroup(num)
	var w sync.WaitGroup
	w.Add(num)
	for i, c := range g.Clients {
		go func(i int, c *client.Client) {
			defer w.Done()
			c.Connect(cometAddr[i%2])
			// c.Auth("auth message")
			c.JoinRoom(1005)
		}(i, c)
	}
	w.Wait()
	log.Printf("%d connetc and auth", num)

	msg := "hello"
	g.Send(&msg)
	go func() {
		for {
			time.Sleep(time.Second)
			sum := 0
			for _, c := range g.Clients {
				sum += c.ReceiveCount
			}
		}
	}()
	select {}
}
