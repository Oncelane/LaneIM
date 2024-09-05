package model_test

import (
	"laneIM/proto/msg"
	"laneIM/src/common"
	"laneIM/src/config"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"testing"
)

// func TestRoom(t *testing.T) {
// 	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
// 	v, err := model.RoomGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}
// 	log.Println(v.String())
// 	if v == nil {
// 		v = &msg.RoomInfo{}
// 	} else {
// 		v.Reset()
// 	}
// 	v.Roomid = 1
// 	v.Server["testHost"] = true
// 	model.RoomSet(r.Client, v)

// 	v, err = model.RoomGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}
// 	log.Println("value:", v.String())
// 	v.Reset()

// 	v, err = model.RoomGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}
// 	if _, exist := v.Server["testHost"]; !exist {
// 		t.Error("set err")
// 	}
// 	log.Println(v.String())
// }

// func TestUser(t *testing.T) {
// 	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7005"})
// 	v, err := model.UserGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}

// 	v.Reset()

// 	v.Userid = 1
// 	v.Roomid[1] = true
// 	v.Server["testHost"] = true
// 	err = model.UserSet(r.Client, v)
// 	if err != nil {
// 		t.Error("set err")
// 	}
// 	v.Reset()

// 	v, err = model.UserGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}
// 	if _, exist := v.Server["testHost"]; !exist {
// 		t.Error("set err")
// 	}
// 	log.Println(v.String())
// }

// 初始化1005房间
// 四个用户进行广播
func TestInitRoomAndUser(t *testing.T) {
	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
	e := pkg.NewEtcd(config.Etcd{Addr: []string{"127.0.0.1:2379"}})
	e.SetAddr("redis:1", "127.0.0.1:7001")
	e.SetAddr("redis:2", "127.0.0.1:7003")
	e.SetAddr("redis:3", "127.0.0.1:7006")
	// usermap := map[int64]bool{21: true, 22: true, 23: true, 24: false}

	userid := []int64{21, 22, 23, 24}
	roomid := []int64{1005, 1005, 1005, 1005}
	online := []bool{true, true, true, false}
	server := []string{"127.0.0.1:50050", "127.0.0.1:50051", "127.0.0.1:50050", "127.0.0.1:50051"}

	user := msg.UserInfo{}
	log.Println("写入user")
	for i := range len(userid) {
		user.Reset()
		user.Roomid = make(map[int64]bool)
		user.Server = make(map[string]bool)
		user.Userid = userid[i]
		user.Roomid[roomid[i]] = true
		user.Online = online[i]
		user.Server[server[i]] = true
		_, err := model.UserDel(r.Client, common.Int64(user.Userid))
		if err != nil {
			t.Error(err)
		}
		err = model.UserNew(r.Client, common.Int64(user.Userid))
		if err != nil {
			t.Error(err)
		}
		if user.Online {
			model.UserOnline(r.Client, common.Int64(user.Userid), server[i])
		}
	}
	log.Println("查看user")
	getUserId, err := model.AllUserid(r.Client)
	if err != nil {
		t.Error(err)
	}
	log.Println(getUserId)

	testroom := msg.RoomInfo{
		Roomid: 1005,
		// Users:     usermap,
		// OnlineNum: 3,
		// Server:    map[string]bool{"127.0.0.1:50050": true},
	}
	_, err = model.RoomDel(r.Client, common.Int64(testroom.Roomid))
	if err != nil {
		t.Error(err)
	}
	log.Println("创建room")
	err = model.RoomNew(r.Client, common.Int64(testroom.Roomid), common.Int64(userid[0]), server[0])
	if err != nil {
		t.Error(err)
	}
	log.Println("给room添加user")
	log.Println("给user添加room")
	for i, id := range userid {
		err = model.RoomJoinUser(r.Client, common.Int64(testroom.Roomid), common.Int64(id))
		if err != nil {
			t.Error(err)
		}
		err = model.RoomPutComet(r.Client, common.Int64(testroom.Roomid), server[i])
		if err != nil {
			t.Error(err)
		}
		err = model.UserJoinRoomid(r.Client, common.Int64(id), common.Int64(testroom.Roomid))
		if err != nil {
			t.Error(err)
		}
	}

	log.Println("查看user")
	getUserId, err = model.RoomQueryUserid(r.Client, common.Int64(testroom.Roomid))
	if err != nil {
		t.Error(err)
	}
	log.Println(getUserId)
	log.Println("查看comet")
	cometAddr, err := model.RoomQueryComet(r.Client, common.Int64(testroom.Roomid))
	if err != nil {
		t.Error(err)
	}
	log.Println(cometAddr)

	log.Println("查询user的roomid")
	userRoomid, err := model.UserQueryRoomid(r.Client, common.Int64(userid[0]))
	if err != nil {
		t.Error(err)
	}
	log.Println(userRoomid)

}
