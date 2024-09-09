package dao_test

import (
	osql "database/sql"
	"errors"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"testing"
	"time"
)

// func TestRoom(t *testing.T) {
// 	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
// 	v, err := dao.RoomGet(r.Client, 1)
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
// 	dao.RoomSet(r.Client, v)

// 	v, err = dao.RoomGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}
// 	log.Println("value:", v.String())
// 	v.Reset()

// 	v, err = dao.RoomGet(r.Client, 1)
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
// 	v, err := dao.UserGet(r.Client, 1)
// 	if err != nil {
// 		t.Error("get err")
// 	}

// 	v.Reset()

// 	v.Userid = 1
// 	v.Roomid[1] = true
// 	v.Server["testHost"] = true
// 	err = dao.UserSet(r.Client, v)
// 	if err != nil {
// 		t.Error("set err")
// 	}
// 	v.Reset()

// 	v, err = dao.UserGet(r.Client, 1)
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

func TestGetUserOnline(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001"}})
	rt, cometAddr, err := dao.UserOnline(rdb.Client, db, 21)
	if err != nil {
		t.Error(err)
	}
	log.Println("query redis and sql", rt, "comet:", cometAddr)
}

func TestSetUserOnline(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	err := sql.SetUserOnline(db, 21, "localhost")
	if err != nil {
		t.Error(err)
	}
}

func TestSetUserOffline(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	err := sql.SetUseroffline(db, 21)
	if err != nil {
		t.Error(err)
	}
}

func TestRedisAndSqlAllRoomid(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	rdb := pkg.NewRedisClient(config.Redis{Addr: []string{"127.0.0.1:7001"}})
	cache := localCache.Cache(time.Second * 30)
	rt, err := dao.RoomUserid(cache, rdb.Client, db, 1005)
	if err != nil {
		t.Error(err)
	}
	log.Println("query redis and sql", rt)
}

func TestAddRoomUser(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	err := sql.AddRoomUser(db, 1006, 21)
	if err != nil {
		t.Error(err)
	}
}

func TestDelRoomUser(t *testing.T) {
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	err := sql.DelRoomUser(db, 1006, 21)
	if err != nil {
		t.Error(err)
	}
}

func TestInitRoomAndUser(t *testing.T) {
	// r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
	e := pkg.NewEtcd(config.Etcd{Addr: []string{"127.0.0.1:2379"}})
	e.SetAddr("redis:1", "127.0.0.1:7001")
	e.SetAddr("redis:2", "127.0.0.1:7003")
	e.SetAddr("redis:3", "127.0.0.1:7006")
	mysqlConfig := config.Mysql{}
	mysqlConfig.Default()
	db := sql.DB(mysqlConfig)
	model.Init(db)
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

		err := sql.DelUser(db, user.Userid)
		if err != nil && !errors.Is(err, osql.ErrNoRows) {
			t.Error(err)
		}
		err = sql.NewUser(db, user.Userid)
		if err != nil {
			t.Error(err)
		}
		err = sql.SetUserOnline(db, user.Userid, server[i])
		if err != nil {
			t.Error(err)
		}
	}
	log.Println("查看user")
	getUserId, err := sql.AllUserid(db)
	if err != nil {
		t.Error(err)
	}
	log.Println(getUserId)

	testroom := msg.RoomInfo{
		Roomid: 1005,
	}
	err = sql.DelRoom(db, testroom.Roomid)
	if err != nil {
		t.Error(err)
	}
	log.Println("创建room")
	err = sql.NewRoom(db, testroom.Roomid, userid[0], server[0])
	if err != nil {
		t.Error(err)
	}
	log.Println("给room添加user")
	log.Println("给user添加room")
	for i, id := range userid {
		err = sql.AddRoomUser(db, testroom.Roomid, id)
		if err != nil {
			t.Error(err)
		}
		err = sql.AddRoomComet(db, testroom.Roomid, server[i])
		if err != nil {
			t.Error(err)
		}
		err = sql.AddUserRoom(db, id, testroom.Roomid)
		if err != nil {
			t.Error(err)
		}
	}

	log.Println("查看user")
	getUserId, err = sql.RoomUserid(db, testroom.Roomid)
	if err != nil {
		t.Error(err)
	}
	log.Println(getUserId)
	log.Println("查看comet")
	cometAddr, err := sql.RoomComet(db, testroom.Roomid)
	if err != nil {
		t.Error(err)
	}
	log.Println(cometAddr)

	log.Println("查询user的roomid")
	userRoomid, err := sql.UserRoom(db, userid[0])
	if err != nil {
		t.Error(err)
	}
	log.Println(userRoomid)

}

// func TestInitRoomAndUser(t *testing.T) {
// 	r := pkg.NewRedisClient([]string{"127.0.0.1:7001", "127.0.0.1:7003", "127.0.0.1:7006"})
// 	e := pkg.NewEtcd(config.Etcd{Addr: []string{"127.0.0.1:2379"}})
// 	e.SetAddr("redis:1", "127.0.0.1:7001")
// 	e.SetAddr("redis:2", "127.0.0.1:7003")
// 	e.SetAddr("redis:3", "127.0.0.1:7006")
// 	// usermap := map[int64]bool{21: true, 22: true, 23: true, 24: false}

// 	userid := []int64{21, 22, 23, 24}
// 	roomid := []int64{1005, 1005, 1005, 1005}
// 	online := []bool{true, true, true, false}
// 	server := []string{"127.0.0.1:50050", "127.0.0.1:50051", "127.0.0.1:50050", "127.0.0.1:50051"}

// 	user := msg.UserInfo{}
// 	log.Println("写入user")
// 	for i := range len(userid) {
// 		user.Reset()
// 		user.Roomid = make(map[int64]bool)
// 		user.Server = make(map[string]bool)
// 		user.Userid = userid[i]
// 		user.Roomid[roomid[i]] = true
// 		user.Online = online[i]
// 		user.Server[server[i]] = true
// 		_, err := rds.UserDel(r.Client, common.Int64(user.Userid))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		err = rds.UserNew(r.Client, common.Int64(user.Userid))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if user.Online {
// 			rds.UserSetOnline(r.Client, common.Int64(user.Userid), server[i])
// 		}
// 	}
// 	log.Println("查看user")
// 	getUserId, err := rds.AllUserid(r.Client)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(getUserId)

// 	testroom := msg.RoomInfo{
// 		Roomid: 1005,
// 		// Users:     usermap,
// 		// OnlineNum: 3,
// 		// Server:    map[string]bool{"127.0.0.1:50050": true},
// 	}
// 	err = rds.DelRoom(r.Client, common.Int64(testroom.Roomid))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println("创建room")
// 	err = rds.NewRoom(r.Client, common.Int64(testroom.Roomid), common.Int64(userid[0]), server[0])
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println("给room添加user")
// 	log.Println("给user添加room")
// 	for i, id := range userid {
// 		err = rds.AddRoomUser(r.Client, common.Int64(testroom.Roomid), common.Int64(id))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		err = rds.AddRoomComet(r.Client, common.Int64(testroom.Roomid), server[i])
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		err = rds.AddUserRoom(r.Client, common.Int64(id), common.Int64(testroom.Roomid))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}

// 	log.Println("查看user")
// 	getUserId, err = rds.RoomUserid(r.Client, common.Int64(testroom.Roomid))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(getUserId)
// 	log.Println("查看comet")
// 	cometAddr, err := rds.RoomComet(r.Client, common.Int64(testroom.Roomid))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(cometAddr)

// 	log.Println("查询user的roomid")
// 	userRoomid, err := rds.UserRoom(r.Client, common.Int64(userid[0]))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(userRoomid)

// }
