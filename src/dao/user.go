package dao

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg/laneLog"
	"laneIM/src/pkg/mergewrite"
	"strconv"

	"github.com/allegro/bigcache"
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

func (d *Dao) UserRoom(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) ([]int64, error) {
	key := "user:room:" + strconv.FormatInt(userid, 36)
	r, err := localCache.UserRoomid(cache, userid)
	if err == nil {
		return r, err
	}
	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		// laneLog.Logger.Infoln("user room query redis")
		rt, err := rds.UserMgrRoom(rdb, userid)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			//laneLog.Logger.Infoln("同步到本地cache", rt)
			localCache.SetUserRoomid(cache, userid, rt)
			return rt, nil
		}
		//laneLog.Logger.Infoln("触发sql查询")
		rt, err = db.UserRoom(userid)
		if err != nil {
			return rt, err
		}

		//laneLog.Logger.Infoln("同步到本地cache", rt)
		localCache.SetUserRoomid(cache, userid, rt)
		//laneLog.Logger.Infoln("同步到redis", rt)
		rds.SetNEUSerMgrRoom(rdb, userid, rt)
		return rt, nil
	})
	if r, ok := rt.([]int64); ok {
		return r, err
	} else {
		return nil, fmt.Errorf("batchwriter faild")
	}
}

func (d *Dao) UserRoomBatch(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userids []int64) ([][]int64, error) {
	rt, nonexists := localCache.UserRoomidBatch(cache, userids)

	var redisUseris []int64
	for i, e := range nonexists {
		if e {
			redisUseris = append(redisUseris, userids[i])
		}
	}
	// laneLog.Logger.Infof("命中本地缓存UserRoomBatch [%d/%d]", len(rt)-len(redisUseris), len(rt))

	if len(redisUseris) == 0 {
		return rt, nil
	}

	redisrt, nonexistsRedis := rds.UserMgrRoomBatch(rdb, redisUseris)
	var sqlUserids []int64

	for i, e := range nonexistsRedis {
		if e {
			sqlUserids = append(sqlUserids, redisUseris[i])
		}
	}
	// laneLog.Logger.Infof("命中redis缓存UserRoomBatch [%d/%d]", len(redisrt)-len(sqlUserids), len(redisrt))
	if len(redisrt)-len(sqlUserids) != 0 {
		var redisIndex = 0
		for i, e := range nonexists {
			if e {
				rt[i] = redisrt[redisIndex]
				nonexists[i] = false
				redisIndex++
			}
		}
		localCache.SetUserRoomidBatch(cache, redisUseris, redisrt)
		if len(sqlUserids) == 0 {
			return rt, nil
		}
	}

	// laneLog.Logger.Errorf("触发sql查询UserRoom [%d条] [%v]", len(sqlUserids), sqlUserids)
	sqlrt, err := db.UserRoomBatch(sqlUserids)
	if err != nil {
		laneLog.Logger.Errorln("faild to sql user room batch", err)
		return rt, err
	}
	var sqlindex = 0
	for i, e := range nonexists {
		if e {
			// laneLog.Logger.Debugf("sql 补全 userid[%d] roomids[%v] ", userids[i], sqlrt[sqlindex])
			rt[i] = sqlrt[sqlindex]
			nonexists[i] = false
			sqlindex++
		}
	}
	localCache.SetUserRoomidBatch(cache, sqlUserids, sqlrt)
	// start := time.Now()
	go rds.SetNEUSerMgrRoomBatch(rdb, sqlUserids, sqlrt)
	// laneLog.Logger.Debugln("go SetNEUSerMgrRoomBatch spand", time.Since(start))
	return rt, nil
}

func (d *Dao) UserComet(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) (string, error) {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	r, err := localCache.UserComet(cache, userid)
	if err == nil {
		// laneLog.Logger.Infoln("命中本地缓存UserComet")
		return r, err
	}
	rt, err := d.msergeWriter.Do(key, func() (any, error) {
		// laneLog.Logger.Infoln("触发redis查询UserComet")
		rt, err := rds.UserMgrComet(rdb, userid)
		if err != nil {
			if err != redis.Nil {
				return rt, err
			}
		} else {
			//laneLog.Logger.Infoln("同步到本地cache", rt)
			localCache.SetUserComet(cache, userid, rt)
			return rt, nil
		}
		// laneLog.Logger.Errorln("触发sql查询UserComet")
		rt, err = db.UserComet(userid)
		if err != nil {
			return rt, err
		}
		// laneLog.Logger.Infoln("UserComet同步到本地cache")
		localCache.SetUserComet(cache, userid, rt)
		//TODO 单独同步到redis
		// laneLog.Logger.Infoln("UserComet同步到redis")
		go rds.SetNEUSerMgrComet(rdb, userid, rt)
		return rt, nil
	})
	if r, ok := rt.(string); ok {
		return r, err
	} else {
		return "", fmt.Errorf("batchwriter faild")
	}
}

// func (d *Dao) UserOnline(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userid int64) (bool, string, error) {
// 	rt, cometAddr, err := rds.UserOnline(rdb, userid)
// 	if err != nil {
// 		if err != redis.Nil {
// 			return false, "", err
// 		}
// 	} else {
// 		return rt, cometAddr, nil
// 	}
// 	////laneLog.Logger.Infoln("触发sql查询")
// 	rt, cometAddr, err = db.UserOnlie(userid)
// 	if err != nil {
// 		return rt, cometAddr, err
// 	}
// 	if rt {
// 		rds.SetNEUserOnline(rdb, userid, cometAddr)
// 	} else {
// 		rds.SetNEUserOffline(rdb, userid, cometAddr)
// 	}

// 	return rt, cometAddr, nil
// }
