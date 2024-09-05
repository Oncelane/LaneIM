package job

import (
	"context"
	"laneIM/proto/comet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CometClient struct {
	addr   string
	client comet.CometClient

	brodcastCh chan *comet.BrodcastReq
	roomCh     chan *comet.RoomReq
	singleCh   chan *comet.SingleReq
	conn       *grpc.ClientConn
	done       chan struct{}
}

func (j *Job) NewComet(addr string) *CometClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dail faild ", err.Error())
		return nil
	}
	client := comet.NewCometClient(conn)
	c := &CometClient{
		addr:       addr,
		client:     client,
		brodcastCh: make(chan *comet.BrodcastReq, 1024),
		roomCh:     make(chan *comet.RoomReq, 1024),
		singleCh:   make(chan *comet.SingleReq, 1024),
		conn:       conn,
		done:       make(chan struct{}, j.conf.CometRoutineSize),
	}
	log.Println("connet to comet:", addr)
	j.mu.Lock()
	j.comets[addr] = c
	j.mu.Unlock()
	for range j.conf.CometRoutineSize {
		go c.HandlerComet()
	}
	return c
}
func (c *CometClient) HandlerComet() {
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
		case <-c.done:
			goto OUT
		}
	}
OUT:
}
