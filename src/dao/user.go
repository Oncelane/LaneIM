package dao

import (
	"log"

	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func AllUserid(rdb *redis.ClusterClient, db *gorm.DB) ([]int64, error) {
	rt, err := rds.AllUserid(rdb)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.AllUserid(db)
	if err != nil {
		return rt, err
	}
	log.Println("同步到redis")
	rds.SetNEAllUserid(rdb, rt)
	return rt, nil
}

func UserRoom(rdb *redis.ClusterClient, db *gorm.DB, userid int64) ([]int64, error) {
	rt, err := rds.UserRoom(rdb, userid)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.UserRoom(db, userid)
	if err != nil {
		return rt, err
	}
	log.Println("同步到redis")
	err = rds.SetNEUserRoom(rdb, userid, rt)
	if err != nil {
		return rt, err
	}
	return rt, nil
}

func UserComet(rdb *redis.ClusterClient, db *gorm.DB, userid int64) (string, error) {
	rt, err := rds.UserComet(rdb, userid)
	if err != nil {
		if err != redis.Nil {
			return rt, err
		}
	} else {
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = sql.UserComet(db, userid)
	if err != nil {
		return rt, err
	}
	rds.SetNEUserComet(rdb, userid, rt)
	return rt, nil
}

func UserOnline(rdb *redis.ClusterClient, db *gorm.DB, userid int64) (bool, string, error) {
	rt, cometAddr, err := rds.UserOnline(rdb, userid)
	if err != nil {
		if err != redis.Nil {
			return false, "", err
		}
	} else {
		return rt, cometAddr, nil
	}
	log.Println("触发sql查询")
	rt, cometAddr, err = sql.UserOnlie(db, userid)
	if err != nil {
		return rt, cometAddr, err
	}
	if rt {
		rds.SetNEUserOnline(rdb, userid, cometAddr)
	} else {
		rds.SetNEUserOffline(rdb, userid, cometAddr)
	}

	return rt, cometAddr, nil
}
