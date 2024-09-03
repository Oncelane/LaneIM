package comet

import (
	"context"
	"laneIM/proto/comet"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"net"
	"time"

	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Logic struct {
	Addr   string
	Client logic.LogicClient
}

// grpc 连接
func NewLogic(addr string) *Logic {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dail faild ", err.Error())
		return nil
	}
	c := logic.NewLogicClient(conn)
	log.Println("connet to logic:", addr)
	return &Logic{
		Addr:   addr,
		Client: c,
	}
}

type Comet struct {
	mu     sync.RWMutex
	userId int64
	etcd   *pkg.EtcdClient
	logics map[string]*Logic
	conf   config.Comet
	grpc   *grpc.Server

	pool    *pkg.MsgPool
	buckets []*Bucket

	chmu     sync.RWMutex
	channels map[int64]*Channel

	funcRout *WsFuncRouter
}

func NewSerivceComet(conf config.Comet) (ret *Comet) {

	ret = &Comet{
		etcd:   pkg.NewEtcd(conf.Etcd),
		logics: make(map[string]*Logic),
		conf:   conf,
		pool:   pkg.NewMsgPool(),

		//bucket
		buckets: make([]*Bucket, conf.BucketSize),

		//func router
		funcRout: NewWsFuncRouter(),

		channels: make(map[int64]*Channel),
	}
	for i := range ret.buckets {
		ret.buckets[i] = NewBucket()
	}

	// watch logic
	go ret.WatchLogic()

	// server grpc
	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		log.Fatalf("error: logic start faild")
	}
	gServer := grpc.NewServer()
	comet.RegisterCometServer(gServer, ret)
	log.Println("Logic serivce is running on port")
	ret.grpc = gServer
	go func() {
		if err := gServer.Serve(lis); err != nil {
			log.Fatalln("failed to serve : ", err.Error())
		}
	}()

	//init func
	// ret.funcRout.Use("newUser", ret.HandleNewUser)
	ret.funcRout.Use("sendRoom", ret.HandleSendRoom)
	ret.funcRout.Use("queryRoom", ret.HandleRoom)
	ret.funcRout.Use("auth", ret.HandleAuth)

	// regieter etcd
	ret.etcd.SetAddr("grpc:comet/"+conf.Name, conf.Addr)
	return ret
}

func (c *Comet) Close() {
	c.etcd.DelAddr("grpc:comet/"+c.conf.Name, c.conf.Addr)
	c.grpc.Stop()
	log.Println("exit comet")
}

func (c *Comet) WatchLogic() {
	for {
		addrs := c.etcd.GetAddr("grpc:logic")
		c.mu.Lock()

		for _, addr := range addrs {
			// connet to comet
			if _, exist := c.logics[addr]; exist {
				continue
			}
			log.Println("发现logic:", addr)
			c.logics[addr] = NewLogic(addr)
		}
		c.mu.Unlock()
		time.Sleep(time.Second)
	}
}

func (c *Comet) GenUserID() (ret int64) {
	c.mu.Lock()
	ret = c.userId
	c.userId++
	c.mu.Unlock()
	return
}

func (c *Comet) pickLogic() *Logic {
	for {
		c.mu.RLock()
		for _, v := range c.logics {
			return v
		}
		c.mu.RUnlock()
		log.Println("暂无发现logic 5秒后再查询")
		time.Sleep(time.Second)
	}
}
func (c *Comet) Bucket(roomid int64) *Bucket {
	log.Println("choos bucket:", int(roomid)%len(c.buckets))
	return c.buckets[int(roomid)%len(c.buckets)]
}

func (c *Comet) LogicBrodcast(message *logic.SendMsgReq, data []byte) {

}

func (c *Comet) LogictRoom(message *logic.SendMsgReq) {
	_, err := c.pickLogic().Client.SendMsg(context.Background(), message)
	if err != nil {
		log.Panicln(err)
	}
}

func (c *Comet) LogicSingle(message *logic.SendMsgReq) {
	_, err := c.pickLogic().Client.SendMsg(context.Background(), message)
	if err != nil {
		log.Panicln(err)
	}
}

func (c *Comet) Single(context.Context, *comet.SingleReq) (*comet.NoResp, error) {
	return nil, nil
}
func (c *Comet) Brodcast(context.Context, *comet.BrodcastReq) (*comet.NoResp, error) {
	return nil, nil
}
func (c *Comet) Room(_ context.Context, in *comet.RoomReq) (*comet.NoResp, error) {
	log.Println("comet send room")
	c.Bucket(in.Roomid).GetRoom(in.Roomid).Send(&msg.Msg{
		Path: "roomMsg",
		Data: in.Data,
	})

	return nil, nil
}
