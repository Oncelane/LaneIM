package dao

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg/mergewrite"
	"strconv"

	"github.com/go-redis/redis"
)

type Dao struct {
	msergeWriter mergewrite.MergeWriter
}

func NewDao(conf config.BatchWriter) *Dao {
	return &Dao{
		msergeWriter: *mergewrite.NewMergeWriter(conf),
	}
}

func (d *Dao) AllUserid(rdb *redis.ClusterClient, db *sql.SqlDB) ([]int64, error) {
	key := "alluser"
	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		rt, err := rds.AllUserid(rdb)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			return rt, nil
		}
		//log.Println("触发sql查询")
		rt, err = db.AllUserid()
		if err != nil {
			return rt, err
		}
		//log.Println("同步到redis")
		rds.SetNEAllUserid(rdb, rt)
		return rt, nil
	})
	if r, ok := rt.([]int64); ok {
		return r, err
	} else {
		return nil, fmt.Errorf("batchwriter faild")
	}

}

func (d *Dao) UserRoom(rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) ([]int64, error) {
	key := "user:room:" + strconv.FormatInt(userid, 36)
	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		rt, err := rds.UserRoom(rdb, userid)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			return rt, nil
		}
		//log.Println("触发sql查询")

		rt, err = db.UserRoom(userid)
		if err != nil {
			return rt, err
		}
		//log.Println("同步到redis")
		err = rds.SetNEUserRoom(rdb, userid, rt)
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

func (d *Dao) UserComet(rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) (string, error) {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		rt, err := rds.UserComet(rdb, userid)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			return rt, nil
		}
		//log.Println("触发sql查询")
		rt, err = db.UserComet(userid)
		if err != nil {
			return rt, err
		}
		rds.SetNEUserComet(rdb, userid, rt)
		return rt, nil
	})
	if r, ok := rt.(string); ok {
		return r, err
	} else {
		return "", fmt.Errorf("batchwriter faild")
	}
}

func (d *Dao) UserOnline(rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) (bool, string, error) {
	rt, cometAddr, err := rds.UserOnline(rdb, userid)
	if err != nil {
		if err != redis.Nil {
			return false, "", err
		}
	} else {
		return rt, cometAddr, nil
	}
	//log.Println("触发sql查询")
	rt, cometAddr, err = db.UserOnlie(userid)
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
