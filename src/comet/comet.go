package comet

import (
	"context"
	"laneIM/proto/comet"
	"laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao/localCache"
	"laneIM/src/pkg"
	"laneIM/src/pkg/batch"
	"laneIM/src/pkg/laneLog"
	"net"
	"time"

	"github.com/allegro/bigcache"

	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Logic struct {
	Addr   string
	Client logic.LogicClient
	Online bool
	conn   *grpc.ClientConn
}

// grpc 连接
func NewLogic(addr string) *Logic {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		laneLog.Logger.Fatalln("[server] Dail faild ", err.Error())
		return nil
	}
	c := logic.NewLogicClient(conn)
	laneLog.Logger.Infoln("[server] connet to logic:", addr)
	return &Logic{
		Addr:   addr,
		Client: c,
		Online: true,
		conn:   conn,
	}
}

func (l *Logic) Close() {
	l.Online = false
	l.Client = nil
	l.conn.Close()
	//laneLog.Logger.Infoln("remove logic", l.Addr)
}

type Comet struct {
	mu     sync.RWMutex
	etcd   *pkg.EtcdClient
	logics map[string]*Logic
	conf   config.Comet
	grpc   *grpc.Server

	pool    *pkg.MsgPool
	buckets []*Bucket

	// chmu sync.RWMutex
	// channels map[int64]*Channel

	funcRout *WsFuncRouter
	cache    *bigcache.BigCache

	msgUUIDGenerator *pkg.UuidGenerator

	//batch
	BatcherNewUser    *batch.BatchArgs[BatchStructNewUser]
	BatcherSetOnline  *batch.BatchArgs[BatchStructSetOnline]
	BatcherSetOffline *batch.BatchArgs[BatchStructSetOffline]
	BatcherJoinRoom   *batch.BatchArgs[BatchStructJoinRoom]
	BatcherSendRoom   *batch.BatchArgs[BatchStructSendRoom]
	BatcherQueryRoom  *batch.BatchArgs[BatchStructQueryRoom]
	BatcherHis        *batch.BatchArgs[BatchStructQueryStoreMsg]
	BatcherLast       *batch.BatchArgs[BatcheStructLast]
}

func NewSerivceComet(conf config.Comet) (ret *Comet) {

	ret = &Comet{
		conf:  conf,
		pool:  pkg.NewMsgPool(),
		cache: localCache.Cache(time.Minute),
		// channels:         make(map[int64]*Channel),
		msgUUIDGenerator: pkg.NewUuidGenerator(int64(conf.Id)),
	}

	ret.InitBatch()

	ret.InitFunc()

	ret.InitBucket()

	// init etcd
	ret.etcd = pkg.NewEtcd(ret.conf.Etcd)
	// watch logic
	go ret.WatchLogic()

	// server grpc
	ret.ServeGrpc()

	// register
	ret.etcd.SetAddr("grpc:comet:"+conf.Name, conf.WindowIP+conf.GrpcPort)

	return ret
}

func (c *Comet) InitBatch() {
	c.BatcherNewUser = batch.NewBatchArgs(10000, time.Millisecond*100, c.doNewUserBatch)
	// 开启goroutine
	c.BatcherNewUser.Start()

	c.BatcherSetOnline = batch.NewBatchArgs(10000, time.Millisecond*100, c.doSetOnlineBatch)
	c.BatcherSetOnline.Start()

	c.BatcherSetOffline = batch.NewBatchArgs(10000, time.Millisecond*100, c.doSetOfflineBatch)
	c.BatcherSetOffline.Start()

	c.BatcherJoinRoom = batch.NewBatchArgs(10000, time.Millisecond*100, c.doJoinRoomBatch)
	c.BatcherJoinRoom.Start()

	c.BatcherSendRoom = batch.NewBatchArgs(3000, time.Millisecond*100, c.doSendRoomBatch)
	c.BatcherSendRoom.Start()

	c.BatcherQueryRoom = batch.NewBatchArgs(10000, time.Millisecond*100, c.doQueryRoomBatch)
	c.BatcherQueryRoom.Start()

	c.BatcherHis = batch.NewBatchArgs(3000, time.Millisecond*100, c.doQueryStoreMsgBatch)
	c.BatcherHis.Start()

	c.BatcherLast = batch.NewBatchArgs(10000, time.Millisecond*100, c.doQueryLast)
	c.BatcherLast.Start()
}

func (c *Comet) InitFunc() {
	//func router
	c.funcRout = NewWsFuncRouter()
	c.funcRout.Use("sendRoom", c.HandleSendRoomBatch)
	c.funcRout.Use("queryRoom", c.HandleQueryRoomBatch)
	c.funcRout.Use("auth", c.HandleAuth)
	c.funcRout.Use("newUser", c.HandleNewUserBatch)
	c.funcRout.Use("newRoom", c.HandleNewRoom)
	c.funcRout.Use("joinRoom", c.HandleJoinRoomBatch)
	c.funcRout.Use("online", c.HandleSetOnlineBatch)
	c.funcRout.Use("offline", c.HandleSetOfflineBatch)
	c.funcRout.Use("his", c.HandleHisMsg)
	c.funcRout.Use("last", c.HandleLast)
	c.funcRout.Use("subon", c.HandleSubscribe)
	c.funcRout.Use("suboff", c.HandleSubscribeOff)
}

func (c *Comet) InitBucket() {
	c.buckets = make([]*Bucket, c.conf.BucketSize)
	for i := range c.buckets {
		c.buckets[i] = NewBucket()
	}
}

func (c *Comet) ServeGrpc() {
	lis, err := net.Listen("tcp", c.conf.UbuntuIP+c.conf.GrpcPort)
	if err != nil {
		laneLog.Logger.Errorln("error: comet start faild", err)
	}
	gServer := grpc.NewServer()
	comet.RegisterCometServer(gServer, c)
	laneLog.Logger.Infoln("[server] comet serivce is running on port")
	c.grpc = gServer
	go func() {
		if err := gServer.Serve(lis); err != nil {
			log.Fatalln("failed to serve : ", err.Error())
		}
	}()
}

func (c *Comet) Close() {
	c.etcd.DelAddr("grpc:comet:"+c.conf.Name, c.conf.WindowIP+c.conf.GrpcPort)
	c.grpc.Stop()
	laneLog.Logger.Infoln("[server] exit comet")
}

func (c *Comet) WatchLogic() {
	c.logics = make(map[string]*Logic)
	for {
		addrs := c.etcd.GetAddr("grpc:logic")
		remoteAddrs := make(map[string]struct{})
		for _, addr := range addrs {
			remoteAddrs[addr] = struct{}{}
			// connet to comet

			// already exist
			if _, exist := c.logics[addr]; exist {
				continue
			}

			// not exist
			//laneLog.Logger.Infoln("etcd discovery logic:", addr)
			c.mu.Lock()
			c.logics[addr] = NewLogic(addr)
			c.mu.Unlock()
		}

		// exist before but now gone
		for addr, client := range c.logics {
			if _, exist := remoteAddrs[addr]; !exist {
				c.mu.Lock()
				delete(c.logics, addr)
				client.Close()
				c.mu.Unlock()
			}
		}

		time.Sleep(time.Second)
	}
}

func (c *Comet) pickLogic() *Logic {
	for {
		c.mu.RLock()
		for _, v := range c.logics {
			c.mu.RUnlock()
			return v
		}
		c.mu.RUnlock()
		// laneLog.Logger.Infoln("non discovery logic")
		time.Sleep(time.Second)
	}
}

func (c *Comet) Bucket(roomid int64) *Bucket {
	return c.buckets[int(roomid)%len(c.buckets)]
}

// delete channel from all room
func (c *Comet) DelChannel(ch *Channel) {
	// delete(c.channels, ch.id)
	ch.ForceClose()
	for i := range c.buckets {
		c.buckets[i].DelChannelAll(ch)
	}
}

func (c *Comet) DelChannelBatch(in []*BatchStructSetOffline) {
	// for i := range in {
	// 	delete(c.channels, in[i].ch.id)
	// }
	for i := range c.buckets {
		in[i].ch.ForceClose()
		c.buckets[i].DelChannelAllBatch(in)
	}
}

func (c *Comet) LogictSendMsgBatch(message *msg.SendMsgBatchReq) error {
	_, err := c.pickLogic().Client.SendMsgBatch(context.Background(), message)
	if err != nil {
		laneLog.Logger.Fatalln("[server]", err)
	}
	return err
}

func (c *Comet) Single(context.Context, *comet.SingleReq) (*comet.NoResp, error) {
	return nil, nil
}

func (c *Comet) Brodcast(context.Context, *comet.BrodcastReq) (*comet.NoResp, error) {
	return nil, nil
}

// func (c *Comet) Room(_ context.Context, in *comet.RoomReq) (*comet.NoResp, error) {
// 	// laneLog.Logger.Infoln("recv from job", in.Roomid)
// 	c.Bucket(in.Roomid).GetRoom(in.Roomid).Send(&msg.Msg{
// 		Path: "roomMsg",
// 		Data: in.Data,
// 	})

// 	return nil, nil
// }

func (c *Comet) SendMsgBatch(_ context.Context, in *msg.SendMsgBatchReq) (*comet.NoResp, error) {
	//消息处理，userid关注哪些room是需要知道的,由调用loigc的quryroom时得知

	//整合消息，以roomid为单位
	roomMsgBatch := make(map[int64]*msg.MsgBatch)
	for i := range in.Msgs {
		// laneLog.Logger.Debugln("! receive roomid", in.Msgs[i].Roomid, "message:", string(in.Msgs[i].Data))
		if _, exist := roomMsgBatch[in.Msgs[i].Roomid]; !exist {
			// laneLog.Logger.Debugln("pass1")
			roomMsgBatch[in.Msgs[i].Roomid] = new(msg.MsgBatch)
		}
		roomMsgBatch[in.Msgs[i].Roomid].Msgs = append(roomMsgBatch[in.Msgs[i].Roomid].Msgs, &msg.Msg{
			Seq:  -1,
			Path: "receive",
			Data: in.Msgs[i].Data,
		})
	}
	// laneLog.Logger.Infoln("recv from job message count", len(in.Msgs))
	for roomid, msg := range roomMsgBatch {

		// laneLog.Logger.Debugln("pass2")
		c.Bucket(roomid).GetRoom(roomid).SendBatch(msg)
	}

	return nil, nil
}
