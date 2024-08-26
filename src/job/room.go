package job

import (
	"laneIM/proto/comet"
	"sync"
)

// 得有一个结构存储 userid->comet
// 应该是查询 redis， 缓存失效再查询数据库

type Room struct {
	online int
	client map[string]*Comet
	rw     sync.RWMutex
	roomid int64
	users  map[int64]struct{}
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) PutComet(addr string, c *Comet) {
	r.rw.Lock()
	r.client[addr] = c
	r.rw.Unlock()
}

func (r *Room) DelComet(addr string, c *Comet) {
	r.rw.Lock()
	delete(r.client, addr)
	r.rw.Unlock()
}

func (r *Room) Push(message *comet.RoomReq) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.online == 0 {
		return
	}
	for _, c := range r.client {
		c.roomCh <- message
	}
}

func (r *Room) PushSingle(message *comet.SingleReq) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.online == 0 {
		return
	}
	for _, c := range r.client {
		c.singleCh <- message
	}
}

func (r *Room) updateOnline(n int) {
	r.rw.Lock()
	r.online = n
	r.rw.Unlock()
}
