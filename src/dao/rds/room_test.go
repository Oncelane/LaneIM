package rds_test

import (
	"laneIM/src/config"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog.go"
	"testing"
)

func TestRoom(t *testing.T) {
	r := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001"}})
	db := sql.NewDB(config.DefaultMysql())
	room := model.RoomMgr{}
	err := db.DB.Preload("Users").Preload("Comets").First(&room, 1833490758133350400).Error
	if err != nil {
		laneLog.Logger.Infoln("faild to sql roomid")
	}
	laneLog.Logger.Infof("sql room:%+v\n", room)
	err = rds.SetNERoomMgr(r.Client, &room)
	if err != nil {
		t.Error(err)
	}

	getRoom, err := rds.RoomMgr(r.Client, room.RoomID)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infof("redis room:%+v\n", getRoom)
}
