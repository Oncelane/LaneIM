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
	rw     sync.RWMutex
}

func (j *Job) Push(message *comet.RoomReq) {
	j.Bucket(message.Roomid).Room(message)
}

// func (r *Room) PushSingle(message *comet.SingleReq) {
// 	if r.info.OnlineNum == 0 {
// 		return
// 	}
// 	for _, c := range r.client {
// 		c.singleCh <- message
// 	}
// }

func (r *Room) UpdateFromRedis(client *pkg.RedisClient) {
	info, err := model.UpdateRoom(client.Client, r.roomid)
	if err != nil {
		log.Println("filed to update room:", r.roomid)
	}
	r.rw.Lock()
	r.info = info
	r.rw.Unlock()
}
