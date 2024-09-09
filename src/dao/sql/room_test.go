package sql_test

import (
	"laneIM/src/config"
	"laneIM/src/dao/sql"
	"testing"
)

func TestAddRoomUseridComet(t *testing.T) {
	s := config.Mysql{}
	s.Default()
	db := sql.DB(s)
	sql.AddRoomCometWithUserid(db, 3, "localhost")
}
