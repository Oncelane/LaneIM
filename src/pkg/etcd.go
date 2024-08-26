package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"laneIM/src/model"
	"log"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	conc "go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdClient struct {
	etcd *clientv3.Client
}

var (
	ErrOpNotAtomic = errors.New("data change not atomic op")
	ErrFaild       = errors.New("op faild for some reason")
)

func NewEtcd() *EtcdClient {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Panicf("faile to connet to etcd: %v", err)
	}
	return &EtcdClient{etcdClient}
}

func (e *EtcdClient) GetAddr(key string) (addrs []string) {
	rt, err := e.etcd.Get(context.Background(), "laneIM/"+key, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return nil
	}
	for _, ev := range rt.Kvs {
		addrs = append(addrs, string(ev.Value))
	}
	return
}

func (e *EtcdClient) SetAddr(key, addr string) {
	_, err := e.etcd.Put(context.Background(), "laneIM/"+key, addr)
	if err != nil {
		log.Fatalln("failed to set etcd:", err.Error())
		return
	}
}

func (e *EtcdClient) GetUserServerAddr(userid int64) ([]string, error) {

	rt, err := e.etcd.Get(context.Background(), UserIdToEtcdServerKey(userid))
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return nil, ErrFaild
	}
	if len(rt.Kvs) == 0 {
		return nil, nil
	}

	return strings.Split(string(rt.Kvs[0].Value), ";"), nil
}

func (e *EtcdClient) SetUserServerAddr(userid int64, addrs []string) error {
	key := UserIdToEtcdServerKey(userid)

	oldAddrs, err := e.GetUserServerAddr(userid)
	if err != nil {
		return err
	}

	oldValue := strings.Join(oldAddrs, ";")
	newValue := strings.Join(addrs, ";")

	err = e.atomicUpdate(key, oldValue, newValue)
	return err
}

func (e *EtcdClient) AddUserServerAddr(userid int64, addr string) error {
	key := UserIdToEtcdServerKey(userid)

	oldAddrs, err := e.GetUserServerAddr(userid)
	if err != nil {
		return err
	}

	for _, v := range oldAddrs {
		if v == addr {
			log.Printf("addr[%s]  already in  userid[%d]  \n", addr, userid)
			return nil
		}
	}

	oldValue := strings.Join(oldAddrs, ";")

	oldAddrs = append(oldAddrs, addr)
	newValue := strings.Join(oldAddrs, ";")
	err = e.atomicUpdate(key, oldValue, newValue)
	return err
}

func (e *EtcdClient) GetUserRoomid(userid int64) (UserRoom, error) {
	rt, err := e.etcd.Get(context.Background(), UserIdToEtcdRoomKey(userid))
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return UserRoom{}, err
	}
	if len(rt.Kvs) == 0 {
		return UserRoom{}, nil
	}
	userRoom := &UserRoom{}
	err = json.Unmarshal(rt.Kvs[0].Value, userRoom)
	return *userRoom, err
}

func (e *EtcdClient) AddUserRoomid(userid int64, roomid int64) error {
	key := UserIdToEtcdRoomKey(userid)

	oldRoom, err := e.GetUserRoomid(userid)
	if err != nil {
		return err
	}
	for _, v := range oldRoom.Id {
		if v == roomid {
			log.Printf("user[%d]  already in  room[%d]  \n", userid, roomid)
			return nil
		}
	}
	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return err
	}

	oldRoom.Id = append(oldRoom.Id, roomid)
	newValue, err := json.Marshal(oldRoom)

	if err != nil {
		return err
	}
	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func (e *EtcdClient) SetUserRoomid(userid int64, roomid []int64) error {
	key := UserIdToEtcdRoomKey(userid)

	oldRoom, err := e.GetUserRoomid(userid)
	if err != nil {
		return err
	}

	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		return err
	}
	newRoom := UserRoom{
		Id: roomid,
	}
	newValue, err := json.Marshal(newRoom)
	if err != nil {
		return err
	}

	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func (e *EtcdClient) DeleteUserRoom(userid int64, roomid int64) error {
	key := UserIdToEtcdRoomKey(userid)

	oldRoom, err := e.GetUserRoomid(userid)
	if err != nil {
		return err
	}
	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return err
	}
	rt := false
	for i, v := range oldRoom.Id {
		if v == roomid {
			oldRoom.Id = append(oldRoom.Id[:i], oldRoom.Id[i+1:]...)
			rt = true
		}
	}

	if !rt {
		log.Printf("user[%d] do not have this room[%d]\n", userid, roomid)
		return nil
	}

	newValue, err := json.Marshal(oldRoom)
	if err != nil {
		return err
	}
	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func UserIdToEtcdRoomKey(userid int64) string {
	return "laneIM/ur/" + strconv.FormatInt(userid, 36)
}

func UserIdToEtcdServerKey(userid int64) string {
	return "laneIM/us/" + strconv.FormatInt(userid, 36)
}

type UserRoom struct {
	Id []int64
}

func (e *EtcdClient) atomicUpdate(key, old, new string) error {
	rt, err := conc.NewSTM(e.etcd, func(s conc.STM) error {
		tmp := s.Get(key)
		log.Printf("key[%v], old[%v], new[%v], tmp[%v]\n", key, old, new, tmp)
		if tmp == old || tmp == "" {
			s.Put(key, new)
			return nil
		}
		return errors.New("update faild, value maybe change")
	})
	if err != nil {
		return err
	}
	if rt.Succeeded {
		return nil
	}
	log.Println("update failed")
	return errors.New("failed")
}

func (e *EtcdClient) GetUserStage(userid int64) (model.UserStage, error) {
	rt, err := e.etcd.Get(context.Background(), UseridToEtcdStageKey(userid))
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return model.UserStage{}, err
	}
	userStage := &model.UserStage{}
	err = json.Unmarshal(rt.Kvs[0].Value, userStage)
	return *userStage, err
}

func (e *EtcdClient) SetUserStage(userid int64, m model.UserStage) error {
	key := UseridToEtcdStageKey(userid)

	oldm, err := e.GetUserStage(userid)
	if err != nil {
		return err
	}

	oldValue, err := json.Marshal(oldm)
	if err != nil {
		return err
	}
	newValue, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func UseridToEtcdStageKey(userid any) string {
	in := userid.(int64)
	return "laneIM/us/" + strconv.FormatInt(in, 36)
}

func (e *EtcdClient) Get(keyI any, genkey func(in any) (out string)) (out []byte, err error) {
	key := genkey(keyI)
	rt, err := e.etcd.Get(context.Background(), key)
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return nil, err
	}
	if len(rt.Kvs) == 0 {
		return nil, nil
	}
	return rt.Kvs[0].Value, nil
}

func (e *EtcdClient) Update(keyI any, genkey func(in any) (out string), f func(in []byte) (out []byte, err error)) error {
	key := genkey(keyI)

	oldValue, err := e.Get(keyI, genkey)
	if err != nil {
		return err
	}
	newValue, err := f(oldValue)
	if err != nil {
		return err
	}

	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func (e *EtcdClient) Set(keyI any, genkey func(in any) (out string), m any) error {
	key := genkey(keyI)

	oldValue, err := e.Get(keyI, genkey)
	if err != nil {
		return err
	}
	newValue, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = e.atomicUpdate(key, string(oldValue), string(newValue))
	return err
}

func (e *EtcdClient) SetUserOnline(userid int64) error {
	return e.Update(userid, UseridToEtcdStageKey, func(in []byte) (out []byte, err error) {
		ustage := &model.UserStage{}
		err = json.Unmarshal(in, ustage)
		if err != nil {
			return nil, err
		}
		ustage.Online = true
		out, err = json.Marshal(ustage)
		return
	})
}

func (e *EtcdClient) SetUserOffline(userid int64) error {
	return e.Update(userid, UseridToEtcdStageKey, func(in []byte) (out []byte, err error) {
		ustage := &model.UserStage{}
		err = json.Unmarshal(in, ustage)
		if err != nil {
			return nil, err
		}
		ustage.Online = false
		out, err = json.Marshal(ustage)
		return
	})
}

func (e *EtcdClient) NewUser(m model.UserStage) error {
	return e.Set(m.Userid, UseridToEtcdStageKey, m)
}
