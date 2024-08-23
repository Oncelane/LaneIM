package comet

import (
	"context"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/pkg"

	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Logic struct {
	Addr   string
	Client logic.LogicClient
}

func NewLogic(addr string) *Logic {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dail faild ", err.Error())
		return nil
	}
	c := logic.NewLogicClient(conn)
	return &Logic{
		Addr:   addr,
		Client: c,
	}
}

type UserID int
type Comet struct {
	mu     *sync.Mutex
	userId UserID
	etcd   *pkg.EtcdClient
	logics []*Logic
}

var (
	logicName string = "logic"
)

func NewComet() (ret *Comet) {

	ret = &Comet{
		etcd: pkg.NewEtcd(),
	}
	addrs := ret.etcd.GetAddr(logicName)
	for _, add := range addrs {
		logic := NewLogic(add)
		if logic == nil {
			log.Fatalf("failt to connet add:%s\n", add)
		}
		ret.logics = append(ret.logics, logic)
	}
	return ret
}

func (c *Comet) GenUserID() (ret UserID) {
	c.mu.Lock()
	ret = c.userId
	c.userId++
	c.mu.Unlock()
	return
}

func (c *Comet) pickLogic() *Logic {
	return c.logics[rand.Int()%len(c.logics)]
}

func (c *Comet) LogicBrodcast(message *msg.Msg) {
	req := logic.BrodcastReq{
		Msg: message,
	}
	_, err := c.pickLogic().Client.Brodcast(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
}

func (c *Comet) LogicBrodcastRoom(message *msg.Msg, roomid int64) {
	req := logic.BrodcastRoomReq{
		Msg:    message,
		Roomid: roomid,
	}
	_, err := c.pickLogic().Client.BrodcastRoom(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
}

func (c *Comet) LogicSingle(message *msg.Msg, userid UserID, roomid int64) {
	req := logic.SingleReq{
		Msg:    message,
		Roomid: roomid,
		Userid: int64(userid),
	}
	_, err := c.pickLogic().Client.Single(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
}
