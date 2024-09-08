package job

import (
	"laneIM/proto/comet"
	"laneIM/proto/msg"
	"laneIM/src/dao"
	"laneIM/src/pkg"
	"log"
	"sync"

	"github.com/allegro/bigcache"
	"gorm.io/gorm"
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

func (r *Room) UpdateFromRedis(cache *bigcache.BigCache, rds *pkg.RedisClient, db *gorm.DB) error {
	serversMap := make(map[string]bool)
	servers, err := dao.RoomComet(cache, rds.Client, db, r.roomid)
	if err != nil {
		log.Printf("faild to read room:%d 's comets\n", err)
		return err
	}
	for _, member := range servers {
		serversMap[member] = true
	}
	useridMap := make(map[int64]bool)
	usersid, err := dao.RoomUserid(cache, rds.Client, db, r.roomid)
	if err != nil {
		log.Printf("faild to read room:%d 's suerids\n", err)
		return err
	}
	for _, member := range usersid {
		// TODO:默认在线状态
		useridMap[member] = true
	}
	r.rw.Lock()
	r.info = &msg.RoomInfo{
		Roomid: r.roomid,
		Server: serversMap,
		Users:  useridMap,
	}
	r.rw.Unlock()
	return nil
}
