package canal_test

import (
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
)

func TestRoomCometNotExist(t *testing.T) {
	conf := config.Canal{}
	config.Init("../cmd/canal/config.yml", &conf)
	roomid := 1837816012817301504
	db := sql.NewDB(conf.Mysql)
	redis := pkg.NewRedisClient(conf.Redis)
	daoo := dao.NewDao(config.DefaultBatchWriter())
	err := daoo.UpdateCacheRoomComet(redis.Client, db, int64(roomid))
	if err != nil {
		t.Fatal(err)
	}
	rtt, err := rds.RoomCometBatch(redis.Client, []int64{int64(roomid)})
	if err != nil {
		t.Fatal(err)
	}
	laneLog.Logger.Infoln("query redis roomid=", roomid, "comet=", rtt)
}

func TestRoomCometExist(t *testing.T) {
	conf := config.Canal{}
	config.Init("../cmd/canal/config.yml", &conf)
	roomid := 1837459066125811712
	// db := sql.NewDB(conf.Mysql)
	redis := pkg.NewRedisClient(conf.Redis)
	// daoo := dao.NewDao(config.DefaultBatchWriter())

	err := rds.SetNERoomComet(redis.Client, int64(roomid), []string{"i am dirty data"})
	if err != nil {
		t.Fatal(err)
	}

	rtt, err := rds.RoomCometBatch(redis.Client, []int64{int64(roomid)})
	if err != nil {
		t.Fatal(err)
	}
	if rtt[0][0] != "i am dirty data" {
		t.Fatal("dirty set wrong")
	}

	// err = rds.SetEXRoomComet(redis.Client, int64(roomid), []string{"i am new data"})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// err = daoo.UpdateCacheRoomComet(redis.Client, db, int64(roomid))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// rtt, err = rds.RoomCometBatch(redis.Client, []int64{int64(roomid)})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// laneLog.Logger.Infoln("query redis roomid=", roomid, "comet=", rtt)
}
