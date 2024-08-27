package comet

import (
	"context"
	"laneIM/proto/logic"
	"laneIM/src/config"
	"laneIM/src/pkg"
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
	return &Logic{
		Addr:   addr,
		Client: c,
	}
}

type Comet struct {
	mu     *sync.RWMutex
	userId int64
	etcd   *pkg.EtcdClient
	logics map[string]*Logic
}

func NewSerivceComet(conf config.Comet) (ret *Comet) {

	ret = &Comet{
		etcd: pkg.NewEtcd(),
	}

	// 注册自己
	ret.etcd.SetAddr("comet/1", "127.0.0.1:50051")

	// 发现logic
	go ret.WatchLogic()
	return ret
}

func (c *Comet) WatchLogic() {
	for {
		addrs := c.etcd.GetAddr("logic")
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
		for _, v := range c.logics {
			return v
		}
		log.Println("暂无发现logic 5秒后再查询")
		time.Sleep(5)
	}
}

func (c *Comet) LogicBrodcast(message *logic.SendMsgReq, data []byte) {

}

func (c *Comet) LogictRoom(message *logic.SendMsgReq, data []byte, userid int64, roomid int64) {
	req := logic.SendMsgReq{
		Data:   data,
		Roomid: roomid,
		Path:   "brodcast",
		Userid: userid,
	}
	_, err := c.pickLogic().Client.SendMsg(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
}

func (c *Comet) LogicSingle(message *logic.SendMsgReq, data []byte, userid int64, touserid int64, roomid int64) {
	req := logic.SendMsgReq{
		Data:     data,
		Path:     "brodcast",
		Userid:   userid,
		ToUserid: userid,
	}
	_, err := c.pickLogic().Client.SendMsg(context.Background(), &req)
	if err != nil {
		log.Panicln(err)
	}
}
