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
	db.AddRoomCometWithUserid(3, "localhost")
}

func TestSlowSql(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.NewDB(s)
	model.Init(db.DB)
	startTime := time.Now()
	// db.NewRoom(1005, 1, "localhost")
	for i := range 1000 {
		log.Println("i", i)
		// err := sql.AddRoomUser(db, 1005, int64(i))
		err := db.NewRoom(int64(i), int64(i), "localhost")
		if err != nil {
			t.Error(err)
		}
	}
	log.Println("spand time:", time.Since(startTime))
}
