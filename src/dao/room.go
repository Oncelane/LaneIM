package dao

import (
	"log"

	"laneIM/src/dao/localCache"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func AllRoomid(rdb *redis.ClusterClient, db *gorm.DB) ([]int64, error) {
	rt, err := rds.AllRoomid(rdb)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.AllRoomid(db)
	if err != nil {
		return rt, err
	}
	if len(rt) == 0 {
		return rt, nil
	}
	log.Println("同步到redis")
	rds.SetNEAllRoomid(rdb, rt)
	return rt, nil
}

func RoomUserid(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *gorm.DB, roomid int64) ([]int64, error) {
	rt, err := localCache.RoomUserid(cache, roomid)
	if err == nil {
		return rt, err
	}
	log.Println("触发redis查询")
	rt, err = rds.RoomUserid(rdb, roomid)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		log.Println("同步到本地cache")
		localCache.SetRoomUserid(cache, roomid, rt)
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.RoomUserid(db, roomid)
	if err != nil {
		return rt, err
	}
	if len(rt) == 0 {
		return rt, nil
	}
	log.Println("同步到redis")

	err = rds.SetNERoomUser(rdb, roomid, rt)
	if err != nil {
		return rt, err
	}
	return rt, nil
}

func RoomComet(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *gorm.DB, roomid int64) ([]string, error) {
	rt, err := localCache.RoomComet(cache, roomid)
	if err == nil {
		return rt, err
	}
	log.Println("触发redis查询")
	rt, err = rds.RoomComet(rdb, roomid)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		log.Println("同步到本地cache")
		localCache.SetRoomComet(cache, roomid, rt)
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.RoomComet(db, roomid)
	if err != nil {
		return rt, err
	}
	if len(rt) == 0 {
		return rt, nil
	}
	log.Println("同步到redis")
	rds.SetNERoomComet(rdb, roomid, rt)
	return rt, nil
}
