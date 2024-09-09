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
	db := sql.DB(s)
	sql.AddRoomCometWithUserid(db, 3, "localhost")
}

func TestSlowSql(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.DB(s)
	model.Init(db)
	startTime := time.Now()
	sql.NewRoom(db, 1005, 1, "localhost")
	for i := range 1000 {
		// err := sql.AddRoomUser(db, 1005, int64(i))
		err := sql.AddRoomUser(db, 1005, int64(i))
		if err != nil {
			t.Error(err)
		}
	}
	log.Println("spand time:", time.Since(startTime))
}
