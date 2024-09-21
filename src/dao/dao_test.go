package dao_test

import (
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis/v8"
)

func TestCahce(t *testing.T) {
	lcache := localCache.Cache(time.Millisecond * 1)
	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"}})
	db := sql.NewDB(config.DefaultMysql())

	dao.Init(db.DB, nil)
	num := 100
	userids := make([]int64, num)
	for i := range num {
		userids[i] = int64(i) + 1834401281997800424
	}
	_, err := UserRoomBatch(lcache, rdb.Client, db, userids)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("query1:")
	time.Sleep(time.Second * 5)
	_, err = UserRoomBatch(lcache, rdb.Client, db, userids)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("query2:")
}

func UserRoomBatch(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userids []int64) ([][]int64, error) {
	rt, nonexists := localCache.UserRoomidBatch(cache, userids)

	var redisUseris []int64
	for i, e := range nonexists {
		if e {
			redisUseris = append(redisUseris, userids[i])
		}
	}
	laneLog.Logger.Infof("命中本地缓存UserRoomBatch [%d/%d]", len(rt)-len(redisUseris), len(rt))

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
	laneLog.Logger.Infof("命中redis缓存UserRoomBatch [%d/%d]", len(redisrt)-len(sqlUserids), len(redisrt))
	if len(redisrt)-len(sqlUserids) != 0 {
		var redisIndex = 0
		for i, e := range nonexists {
			if !e {
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

	laneLog.Logger.Errorf("触发sql查询UserRoom [%d条]", len(sqlUserids))
	sqlrt, err := db.UserRoomBatch(sqlUserids)
	if err != nil {
		laneLog.Logger.Errorln("faild to sql user room batch", err)
		return rt, err
	}
	var sqlindex = 0
	for i, e := range nonexists {
		if !e {
			// laneLog.Logger.Debugln("sql 补全", i)
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
