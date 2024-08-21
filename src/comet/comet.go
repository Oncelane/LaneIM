package comet

import (
	"context"
	"laneIM/proto/pb"
	"laneIM/src/pkg"
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Logic struct {
	Addr string
	C    pb.ExampleServiceClient
}

func NewLogic(addr string) *Logic {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dail faild ", err.Error())
		return nil
	}
	c := pb.NewExampleServiceClient(conn)
	return &Logic{
		Addr: addr,
		C:    c,
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

func (c *Comet) Brodcast(msg string) {
	req := pb.HelloReq{
		Name: msg,
	}
	resp, err := c.pickLogic().C.SayHello(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(resp.Message)
}
