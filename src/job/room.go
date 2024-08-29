package job

import (
	"laneIM/proto/comet"
	"laneIM/proto/msg"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"sync"
)

// 得有一个结构存储 userid->comet
// 应该是查询 redis， 缓存失效再查询数据库

type Room struct {
	roomid int64
	info   *msg.RoomInfo
	client map[string]*ClientComet
	rw     sync.RWMutex
}

func NewRoom(id int64) *Room {
	return &Room{
		roomid: id,
		client: make(map[string]*ClientComet),
	}
}

func (r *Room) PutComet(addr string, c *ClientComet) {
	r.rw.Lock()
	r.client[addr] = c
	r.rw.Unlock()
}

func (r *Room) DelComet(addr string, c *ClientComet) {
	r.rw.Lock()
	delete(r.client, addr)
	r.rw.Unlock()
}

func (r *Room) Push(message *comet.RoomReq) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.info.OnlineNum == 0 {
		return
	}
	for _, c := range r.client {
		c.roomCh <- message
	}
}

func (r *Room) PushSingle(message *comet.SingleReq) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.info.OnlineNum == 0 {
		return
	}
	for _, c := range r.client {
		c.singleCh <- message
	}
}

func (r *Room) UpdateFromRedis(client *pkg.RedisClient) {
	info, err := model.RoomGet(client.Client, r.roomid)
	if err != nil {
		log.Println("filed to update room:", r.roomid)
	}
	r.rw.Lock()
	r.info = info
	r.rw.Unlock()
}
