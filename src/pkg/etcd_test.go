package pkg_test

import (
	"encoding/json"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"reflect"
	"strings"
	"testing"
)

var etcd pkg.EtcdClient

func init() {
	conf := config.Etcd{
		Addr: []string{"127.0.0.1:2379"},
	}
	etcd = *pkg.NewEtcd(conf)
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestEtcd(t *testing.T) {
	roomid := []int64{4, 2, 1, 3, 4}
	newRoom := pkg.UserRoom{
		Id: roomid,
	}
	newValue, _ := json.Marshal(newRoom)

	err := etcd.SetUserRoomid(1, roomid)
	if err != nil {
		t.Error(err)
	}
	a, err := etcd.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	setValue, _ := json.Marshal(a)
	if !reflect.DeepEqual(setValue, newValue) {
		t.Errorf("wrong set")
	}
	newRoom.Id = append(roomid, 2345)
	newValue, _ = json.Marshal(newRoom)
	laneLog.Logger.Infoln(etcd.AddUserRoomid(1, 2345))
	a, err = etcd.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	setValue, _ = json.Marshal(a)
	if !reflect.DeepEqual(setValue, newValue) {
		t.Errorf("wrong add")
	}
	v, err := etcd.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln(v)
}

func TestEtcdServerAddr(t *testing.T) {
	addrs := []string{"1.1.1.1", "2.2.2.2"}
	newValue := strings.Join(addrs, ";")

	err := etcd.SetUserServerAddr(1, addrs)
	if err != nil {
		t.Error(err)
	}
	setAddrs, err := etcd.GetUserServerAddr(1)
	setValue := strings.Join(setAddrs, ";")
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln("setValue:", setValue)
	if !reflect.DeepEqual(newValue, setValue) {
		t.Errorf("wrong set")
	}

	laneLog.Logger.Infoln(etcd.AddUserServerAddr(1, "3.3.3.3"))
	addrs = append(addrs, "3.3.3.3")
	newValue = strings.Join(addrs, ";")
	setAddrs, err = etcd.GetUserServerAddr(1)
	if err != nil {
		t.Error(err)
	}
	setValue = strings.Join(setAddrs, ";")
	if !reflect.DeepEqual(newValue, setValue) {
		t.Errorf("wrong add")
	}
	value, err := etcd.GetUserServerAddr(1)
	if err != nil {
		t.Error(err)
	}
	laneLog.Logger.Infoln(value)
}

// func TestEtcdSetWithFunc(t *testing.T) {
// 	testUser := model.UserStage{
// 		Userid:  1,
// 		Online:  false,
// 		Server:  "localhost",
// 	}
// 	err := etcd.NewUser(testUser)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	v, err := etcd.Get(testUser.Userid, pkg.UseridToEtcdStageKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	getUser := &model.UserStage{}
// 	json.Unmarshal(v, getUser)
// 	laneLog.Logger.Infoln(getUser)

// 	etcd.SetUserOnline(testUser.Userid)

// 	v, err = etcd.Get(testUser.Userid, pkg.UseridToEtcdStageKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	getUser = &model.UserStage{}
// 	json.Unmarshal(v, getUser)
// 	laneLog.Logger.Infoln(getUser)
// }
