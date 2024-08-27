package job

import (
	"context"
	"laneIM/proto/comet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientComet struct {
	addr       string
	client     comet.CometClient
	brodcastCh chan *comet.BrodcastReq
	roomCh     chan *comet.RoomReq
	singleCh   chan *comet.SingleReq
}

func NewComet(addr string) *ClientComet {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dail faild ", err.Error())
		return nil
	}
	client := comet.NewCometClient(conn)
	c := &ClientComet{
		addr:       addr,
		client:     client,
		brodcastCh: make(chan *comet.BrodcastReq, 1024),
		roomCh:     make(chan *comet.RoomReq, 1024),
		singleCh:   make(chan *comet.SingleReq, 1024),
	}
	go c.HandlerComet()
	return c
}

func (c *ClientComet) HandlerComet() {
	for {
		select {
		case msg := <-c.brodcastCh:
			_, err := c.client.Brodcast(context.Background(), msg)
			if err != nil {
				log.Println("brodcrast err:", err)
				continue
			}
		case msg := <-c.roomCh:
			_, err := c.client.Room(context.Background(), msg)
			if err != nil {
				log.Println("brodcrastRoom err:", err)
				continue
			}

		case msg := <-c.singleCh:
			_, err := c.client.Single(context.Background(), msg)
			if err != nil {
				log.Println("single err:", err)
				continue
			}
		}
	}
}
