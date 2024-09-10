package sql_test

import (
	"laneIM/src/config"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"log"
	"testing"
	"time"
)

func TestAddRoomUseridComet(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.NewDB(s)
	log.Println(db.RoomComet(1833434241438842880))
}

var addr = []string{"127.0.0.1:50051", "127.0.0.1:50050"}

func TestSlowSql(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.NewDB(s)
	model.Init(db.DB)
	startTime := time.Now()
	db.NewUser(1)
	db.NewRoom(1005, 1, "localhost")
	for i := range 999 {
		log.Println("i", i+1)
		db.NewUser(int64(i) + 1)
		err := db.AddRoomUser(1005, int64(i)+1)
		db.SetUserOnline(int64(i)+1, addr[i%2])
		db.AddRoomCometWithUserid(int64(i)+1, addr[i%2])
		// err := db.NewRoom(int64(i), int64(i), "localhost")
		if err != nil {
			t.Error(err)
		}
	}
	log.Println("spand time:", time.Since(startTime))
}
