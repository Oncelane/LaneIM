package sql_test

import (
	"laneIM/src/config"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"laneIM/src/pkg/laneLog.go"
	"testing"
	"time"
)

func TestAddRoomUseridComet(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.NewDB(s)
	laneLog.Logger.Infoln(db.RoomComet(1833469140778614784))
}

var addr = []string{"127.0.0.1:50051", "127.0.0.1:50050"}
var num = 3999

func TestSlowSql(t *testing.T) {

	db := sql.NewDB(config.DefaultMysql())
	model.Init(db.DB)
	startTime := time.Now()
	db.NewUser(1)
	db.NewRoom(1005, 1, "localhost")
	var newUsertime time.Duration
	var addroomUsertime time.Duration
	var setUseronlinetime time.Duration
	var addroomcomettime time.Duration
	for i := range num {
		laneLog.Logger.Infoln("i", i+1)
		start := time.Now()
		db.NewUser(int64(i) + 1)
		newUsertime += time.Since(start)

		start = time.Now()
		err := db.AddRoomUser(1005, int64(i)+1)
		addroomUsertime += time.Since(start)

		start = time.Now()
		db.SetUserOnline(int64(i)+1, addr[i%2])
		setUseronlinetime += time.Since(start)
		start = time.Now()
		db.AddRoomCometWithUserid(int64(i)+1, addr[i%2])
		addroomcomettime += time.Since(start)
		// err := db.NewRoom(int64(i), int64(i), "localhost")
		if err != nil {
			t.Error(err)
		}
	}
	laneLog.Logger.Infoln("newUsertime spand", newUsertime)
	laneLog.Logger.Infoln("addroomUsertime spand", addroomUsertime)
	laneLog.Logger.Infoln("setUseronlinetime spand", setUseronlinetime)
	laneLog.Logger.Infoln("addroomcomettime", addroomcomettime)
	laneLog.Logger.Infoln("spand time:", time.Since(startTime))
}

func TestSlowSqlOptimines(t *testing.T) {

	db := sql.NewDB(config.DefaultMysql())
	model.Init(db.DB)
	startTime := time.Now()

	db.NewUser(1)
	db.NewRoom(1005, 1, "localhost")
	var newUsertime time.Duration
	var addroomUsertime time.Duration
	var setUseronlinetime time.Duration

	userids := make([]int64, num)

	rooms := make([]int64, num)
	for i := range num {
		userids[i] = int64(i) + 1
		rooms[i] = 1005
	}

	tx := db.DB.Begin()

	start := time.Now()
	db.NewUserBatch(tx, userids)
	newUsertime += time.Since(start)

	start = time.Now()
	db.SetUserOnlineBatch(tx, userids, addr[0])
	setUseronlinetime += time.Since(start)

	start = time.Now()
	db.AddRoomUserBatch(tx, rooms, userids)
	addroomUsertime += time.Since(start)

	laneLog.Logger.Infoln("newUsertime-batch spand", newUsertime)
	laneLog.Logger.Infoln("setUseronlinetime spand", setUseronlinetime)
	laneLog.Logger.Infoln("addroomUsertime spand", addroomUsertime)
	laneLog.Logger.Infoln("spand time:", time.Since(startTime))
}
