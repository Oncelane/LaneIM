package dao_test

import (
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
	"time"
)

// func TestGetUserOnline(t *testing.T) {
// 	mysqlConfig := config.Mysql{}
// 	mysqlConfig.Default()
// 	db := sql.NewDB(mysqlConfig)
// 	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001"}})
// 	d := dao.NewDao(config.DefaultBatchWriter())
// 	rt, cometAddr, err := d.UserOnline(rdb.Client, db, 21)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	laneLog.Logger.Infoln("query redis and sql", rt, "comet:", cometAddr)
// }

func TestSetUserOnline(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.NewDB(mysqlConfig)
	err := db.SetUserOnline(21, "localhost")
	if err != nil {
		t.Error(err)
	}
}

func TestSetUserOffline(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.NewDB(mysqlConfig)
	err := db.SetUseroffline(21)
	if err != nil {
		t.Error(err)
	}
}

func TestRedisAndSqlAllRoomid(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.NewDB(mysqlConfig)
	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001"}})
	cache := localCache.NewLocalCache(time.Minute)
	d := dao.NewDao(config.DefaultBatchWriter())
	rt, err := d.RoomUserid(cache, rdb.Client, db, 1005)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("query redis and sql", rt)
}

func TestAddRoomUser(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.NewDB(mysqlConfig)
	err := db.AddRoomUser(1006, 21)
	if err != nil {
		t.Error(err)
	}
}

func TestDelRoomUser(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.NewDB(mysqlConfig)
	err := db.DelRoomUser(1006, 21)
	if err != nil {
		t.Error(err)
	}
}
