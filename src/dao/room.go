package dao

import (
	"context"
	"fmt"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg/laneLog"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis/v8"
)

// func (d *Dao) AllRoomid(rdb *redis.ClusterClient, db *sql.SqlDB) ([]int64, error) {
// 	rt, err := d.msergeWriter.Do("allroom", func() (any, error) {
// 		rt, err := rds.AllRoomid(rdb)
// 		if err != nil {
// 			if err != redis.Nil {
// 				return rt, err
// 			}
// 		} else {
// 			return rt, nil
// 		}
// 		// //laneLog.Logger.Infoln("触发sql查询")
// 		rt, err = db.AllRoomidSingleflight()
// 		if err != nil {
// 			return rt, err
// 		}
// 		if len(rt) == 0 {
// 			return rt, nil
// 		}
// 		// //laneLog.Logger.Infoln("同步到redis")
// 		rds.SetNEAllRoomid(rdb, rt)
// 		return rt, nil
// 	})
// 	if r, ok := rt.([]int64); ok {
// 		return r, err
// 	} else {
// 		return nil, fmt.Errorf("batchwriter faild")
// 	}

// }

func (d *Dao) RoomUserid(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, roomid int64) ([]int64, error) {

	key := "room:userid" + strconv.FormatInt(roomid, 36)
	r, err := localCache.RoomUserid(cache, roomid)
	if err == nil {
		return r, err
	}
	rt, err := d.msergeWriter.Do(key, func() (any, error) {

		laneLog.Logger.Infoln("[dao] room user query redis")
		rt, err := rds.RoomMgrUserid(rdb, roomid)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			//laneLog.Logger.Infoln("同步到本地cache", rt)
			localCache.SetRoomUserid(cache, roomid, rt)
			return rt, nil
		}
		//laneLog.Logger.Infoln("触发sql查询")
		rt, err = db.RoomUserid(roomid)
		if err != nil {
			return rt, err
		}

		//laneLog.Logger.Infoln("同步到本地cache", rt)
		localCache.SetRoomUserid(cache, roomid, rt)

		// TODO 单独同步某一项到redis
		//laneLog.Logger.Infoln("同步到redis", rt)
		err = rds.SetNERoomMgrUsers(rdb, roomid, rt)
		if err != nil {
			return rt, err
		}

		return rt, nil
	})
	if r, ok := rt.([]int64); ok {
		return r, err
	} else {
		return nil, fmt.Errorf("batchwriter faild")
	}
}

func (d *Dao) RoomComet(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, roomid int64) ([]string, error) {
	startTime := time.Now()
	key := "room:comet" + strconv.FormatInt(roomid, 36)
	r, err := localCache.RoomComet(cache, roomid)
	if err == nil {
		// laneLog.Logger.Debugln("RoomComet命中localcache")
		// laneLog.Logger.Debugln("time on localcache user room spand ", time.Since(startTime))
		return r, err
	}

	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		rtt, err := rds.RoomCometBatch(rdb, []int64{roomid})
		if err != nil {
			if err != redis.Nil {
				return rtt[0], err
			}
		} else {
			laneLog.Logger.Debugln("RoomComet命中redis")
			laneLog.Logger.Debugln("time on redis user room spand ", time.Since(startTime))
			//laneLog.Logger.Infoln("同步到本地cache:", rt)
			localCache.SetRoomComet(cache, roomid, rtt[0])
			return rtt[0], nil
		}
		//laneLog.Logger.Infoln("触发sql查询")
		rt, err := db.RoomComet(roomid)
		if err != nil {
			return rt, err
		}

		//laneLog.Logger.Infoln("同步到本地cache", rt)
		localCache.SetRoomComet(cache, roomid, rt)
		//TODO 单独同步到redis
		//laneLog.Logger.Infoln("同步到redis", rt)
		rds.SetNERoomComet(rdb, roomid, rt)
		return rt, nil
	})
	laneLog.Logger.Debugln("time on sql user room spand ", time.Since(startTime))
	if r, ok := rt.([]string); ok {

		return r, err
	} else {
		return nil, fmt.Errorf("batchwriter faild")
	}

}

func (d *Dao) UpdateCacheRoomComet(rdb *redis.ClusterClient, db *sql.SqlDB, roomid int64) error {
	rrt, err := db.RoomCometBatch([]int64{roomid})
	if err != nil {
		laneLog.Logger.Fatalln(err)
		return err
	}
	comets := rrt[0]
	err = rds.SetEXRoomComet(rdb, roomid, comets)
	if err != nil {
		laneLog.Logger.Fatalln(err)
		return err
	}
	// TODO sub/pub localcache
	_, err = rdb.Publish(context.Background(), "room:comet", roomid).Result()
	if err != nil {
		laneLog.Logger.Fatalln(err)
	}
	return nil
}
