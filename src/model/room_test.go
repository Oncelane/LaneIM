package model_test

import (
	"laneIM/proto/msg"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"testing"
)

func TestRoom(t *testing.T) {
	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
	v, err := model.RoomGet(r.Client, 1)
	if err != nil {
		t.Error("get err")
	}
	log.Println(v.String())
	if v == nil {
		v = &msg.RoomInfo{}
	} else {
		v.Reset()
	}
	v.Roomid = 1
	v.Server = append(v.Server, "testHost")
	model.RoomSet(r.Client, v)

	v, err = model.RoomGet(r.Client, 1)
	if err != nil {
		t.Error("get err")
	}
	log.Println("value:", v.String())
	v.Reset()

	v, err = model.RoomGet(r.Client, 1)
	if err != nil {
		t.Error("get err")
	}
	if v.Server[0] != "testHost" {
		t.Error("set err")
	}
	log.Println(v.String())
}

func TestUser(t *testing.T) {
	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7005"})
	v, err := model.UserGet(r.Client, 1)
	if err != nil {
		t.Error("get err")
	}

	v.Reset()

	v.Userid = 1
	v.Roomid = append(v.Roomid, 1)
	v.Server = append(v.Server, "testHost")
	err = model.UserSet(r.Client, v)
	if err != nil {
		t.Error("set err")
	}
	v.Reset()

	v, err = model.UserGet(r.Client, 1)
	if err != nil {
		t.Error("get err")
	}
	if v.Server[0] != "testHost" {
		t.Error("set err")
	}
	log.Println(v.String())
}

// 初始化1005房间
// 四个用户进行广播
func TestInitRoomAndUser(t *testing.T) {
	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
	e := pkg.NewEtcd()
	e.SetAddr("redis/1", "127.0.0.1:7001")
	e.SetAddr("redis/2", "127.0.0.1:7003")
	e.SetAddr("redis/3", "127.0.0.1:7006")
	usermap := map[int64]bool{21: true, 22: true, 23: true, 24: false}
	testroom := msg.RoomInfo{
		Roomid:    1005,
		Users:     usermap,
		OnlineNum: 3,
	}
	err := model.RoomSet(r.Client, &testroom)
	if err != nil {
		t.Error(err)
	}
	userid := []int64{21, 22, 23, 24}
	roomid := []int64{1005, 1005, 1005, 1005}
	online := []bool{true, true, true, false}
	machine := []int32{0, 0, 0, 0}
	server := []string{"127.0.0.1:50051", "127.0.0.1:50051", "127.0.0.1:50051", "127.0.0.1:50051"}

	user := msg.UserInfo{}
	for i := range len(userid) {
		user.Reset()
		user.Userid = userid[i]
		user.Roomid = append(user.Roomid, roomid[i])
		user.Online = online[i]
		user.Machine = machine[i]
		user.Server = append(user.Server, server[i])

		err := model.UserSet(r.Client, &user)
		if err != nil {
			t.Error(err)
		}
	}

}
