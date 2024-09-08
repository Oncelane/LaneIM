package dao

import (
	"log"

	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func AllRoomid(rdb *redis.ClusterClient, db *gorm.DB) ([]int64, error) {
	rt, err := rds.AllRoomid(rdb)
	if err != nil {
		return rt, err
	}
	if len(rt) != 0 {
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

func RoomUserid(rdb *redis.ClusterClient, db *gorm.DB, roomid int64) ([]int64, error) {
	rt, err := rds.RoomUserid(rdb, roomid)
	if err != nil {
		return rt, err
	}
	if len(rt) != 0 {
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

func RoomComet(rdb *redis.ClusterClient, db *gorm.DB, roomid int64) ([]string, error) {
	rt, err := rds.RoomComet(rdb, roomid)
	if err != nil {
		return rt, err
	}
	if len(rt) != 0 {
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
