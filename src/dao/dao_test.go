package dao_test

import (
	"laneIM/src/config"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog.go"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis"
)

func TestCahce(t *testing.T) {
	lcache := localCache.Cache(time.Millisecond * 1)
	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"}})
	db := sql.NewDB(config.DefaultMysql())
	model.Init(db.DB)
	num := 10
	userids := make([]int64, num)
	for i := range num {
		userids[i] = int64(i) + 1834183641194823680
	}
	rt, err := UserRoomBatch(lcache, rdb.Client, db, userids)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("query1:", rt)
	time.Sleep(time.Second * 5)
	rt, err = UserRoomBatch(lcache, rdb.Client, db, userids)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("query2:", rt)
}

func UserRoomBatch(cache *bigcache.BigCache, rdb *redis.ClusterClient, db *sql.SqlDB, userids []int64) ([][]int64, error) {
	rt, full := localCache.UserRoomidBatch(cache, userids)
	if full {
		laneLog.Logger.Infoln("命中本地缓存UserRoomBatch")
		return rt, nil
	}

	laneLog.Logger.Infoln("UserComet")
	rt, full = rds.UserMgrRoomBatch(rdb, userids)
	if full {
		laneLog.Logger.Infoln("命中redis缓存UserRoomBatch")
		go localCache.SetUserRoomidBatch(cache, userids, rt)
		return rt, nil
	}

	laneLog.Logger.Errorln("redis查询失败： ", rt)
	laneLog.Logger.Errorln("触发sql查询UserRoom")
	rt, err := db.UserRoomBatch(userids)
	if err != nil {
		laneLog.Logger.Errorln("faild to sql user room batch", err)
		return rt, err
	}
	// localCache.SetUserRoomidBatch(cache, userids, rt)
	// start := time.Now()
	rds.SetNEUSerMgrRoomBatch(rdb, userids, rt)
	// laneLog.Logger.Debugln("go SetNEUSerMgrRoomBatch spand", time.Since(start))
	return rt, nil
}
