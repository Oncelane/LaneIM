package job

import (
	"context"
	"laneIM/proto/comet"
	"laneIM/proto/msg"
	"laneIM/src/pkg/batch"
	"laneIM/src/pkg/laneLog.go"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CometClient struct {
	addr   string
	client comet.CometClient

	brodcastCh chan *comet.BrodcastReq
	// roomCh     chan *comet.RoomReq
	roomBatchCh chan *msg.SendMsgBatchReq
	singleCh    chan *comet.SingleReq
	conn        *grpc.ClientConn
	done        chan struct{}

	BatcherSendRoomMsg *batch.BatchArgs[BatchStructSendRoomMsg]
}

func (j *Job) NewComet(addr string) *CometClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		laneLog.Logger.Infoln("Dail faild ", err.Error())
		return nil
	}
	client := comet.NewCometClient(conn)
	c := &CometClient{
		addr:       addr,
		client:     client,
		brodcastCh: make(chan *comet.BrodcastReq, 1024),
		// roomCh:     make(chan *comet.RoomReq, 1024),
		roomBatchCh: make(chan *msg.SendMsgBatchReq, 1024),
		singleCh:    make(chan *comet.SingleReq, 1024),
		conn:        conn,
		done:        make(chan struct{}, j.conf.CometRoutineSize),
	}
	laneLog.Logger.Infoln("connet to comet:", addr)
	j.mu.Lock()
	j.comets[addr] = c
	j.mu.Unlock()

	// io goroutine
	for range j.conf.CometRoutineSize {
		go c.HandlerComet()
	}

	// 启动batcher
	c.BatcherSendRoomMsg = batch.NewBatchArgs(1000, time.Millisecond*100, c.generaDoPushComet())
	c.BatcherSendRoomMsg.Start()
	return c
}

func (c *CometClient) generaDoPushComet() func([]*BatchStructSendRoomMsg) {
	return func(in []*BatchStructSendRoomMsg) {
		// 总消息个数
		var count int
		for i := range in {
			count += len(in[i].arg.Msgs)
		}
		roomBatchSum := make([]*msg.SendMsgReq, count)
		var index int
		for i := range in {
			for j := range in[i].arg.Msgs {
				roomBatchSum[index] = in[i].arg.Msgs[j]
				index++
			}
		}
		// 整合为一个来发送
		c.roomBatchCh <- &msg.SendMsgBatchReq{Msgs: roomBatchSum}
	}
}

func (c *CometClient) HandlerComet() {
	for {
		select {
		case msg := <-c.brodcastCh:
			_, err := c.client.Brodcast(context.Background(), msg)
			if err != nil {
				laneLog.Logger.Infoln("brodcrast err:", err)
				continue
			}
		// case msg := <-c.roomCh:
		// 	_, err := c.client.Room(context.Background(), msg)
		// 	if err != nil {
		// 		laneLog.Logger.Infoln("brodcrastRoom err:", err)
		// 		continue
		// 	}
		case msg := <-c.roomBatchCh:
			_, err := c.client.SendMsgBatch(context.Background(), msg)
			if err != nil {
				laneLog.Logger.Infoln("brodcrastRoom err:", err)
				continue
			}

		case msg := <-c.singleCh:
			_, err := c.client.Single(context.Background(), msg)
			if err != nil {
				laneLog.Logger.Infoln("single err:", err)
				continue
			}
		case <-c.done:
			goto OUT
		}
	}
OUT:
}
