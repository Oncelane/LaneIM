package pkg_test

import (
	"encoding/json"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"reflect"
	"strings"
	"testing"
)

var etcd pkg.EtcdClient

func init() {
	etcd = *pkg.NewEtcd()
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestEtcd(t *testing.T) {
	e := pkg.NewEtcd()

	roomid := []int64{4, 2, 1, 3, 4}
	newRoom := pkg.UserRoom{
		Id: roomid,
	}
	newValue, _ := json.Marshal(newRoom)

	err := e.SetUserRoomid(1, roomid)
	if err != nil {
		t.Error(err)
	}
	a, err := e.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	setValue, _ := json.Marshal(a)
	if !reflect.DeepEqual(setValue, newValue) {
		t.Errorf("wrong set")
	}
	newRoom.Id = append(roomid, 2345)
	newValue, _ = json.Marshal(newRoom)
	log.Println(e.AddUserRoomid(1, 2345))
	a, err = e.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	setValue, _ = json.Marshal(a)
	if !reflect.DeepEqual(setValue, newValue) {
		t.Errorf("wrong add")
	}
	v, err := e.GetUserRoomid(1)
	if err != nil {
		t.Error(err)
	}
	log.Println(v)
}

func TestEtcdServerAddr(t *testing.T) {
	e := pkg.NewEtcd()
	addrs := []string{"1.1.1.1", "2.2.2.2"}
	newValue := strings.Join(addrs, ";")

	err := e.SetUserServerAddr(1, addrs)
	if err != nil {
		t.Error(err)
	}
	setAddrs, err := e.GetUserServerAddr(1)
	setValue := strings.Join(setAddrs, ";")
	if err != nil {
		t.Error(err)
	}
	log.Println("setValue:", setValue)
	if !reflect.DeepEqual(newValue, setValue) {
		t.Errorf("wrong set")
	}

	log.Println(e.AddUserServerAddr(1, "3.3.3.3"))
	addrs = append(addrs, "3.3.3.3")
	newValue = strings.Join(addrs, ";")
	setAddrs, err = e.GetUserServerAddr(1)
	if err != nil {
		t.Error(err)
	}
	setValue = strings.Join(setAddrs, ";")
	if !reflect.DeepEqual(newValue, setValue) {
		t.Errorf("wrong add")
	}
	value, err := e.GetUserServerAddr(1)
	if err != nil {
		t.Error(err)
	}
	log.Println(value)
}

func TestEtcdSetWithFunc(t *testing.T) {
	testUser := model.UserStage{
		Userid:  1,
		Online:  false,
		Server:  "localhost",
		Machine: 1,
	}
	err := etcd.NewUser(testUser)
	if err != nil {
		t.Error(err)
	}
	v, err := etcd.Get(testUser.Userid, pkg.UseridToEtcdStageKey)
	if err != nil {
		t.Error(err)
	}
	getUser := &model.UserStage{}
	json.Unmarshal(v, getUser)
	log.Println(getUser)

	etcd.SetUserOnline(testUser.Userid)

	v, err = etcd.Get(testUser.Userid, pkg.UseridToEtcdStageKey)
	if err != nil {
		t.Error(err)
	}
	getUser = &model.UserStage{}
	json.Unmarshal(v, getUser)
	log.Println(getUser)
}
