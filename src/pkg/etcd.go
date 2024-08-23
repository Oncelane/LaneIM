package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	conc "go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdClient struct {
	*clientv3.Client
}

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

func (e *EtcdClient) GetAddr(key string) (ret []string) {
	rt, err := e.Get(context.Background(), "laneIM/"+key, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return nil
	}
	for _, ev := range rt.Kvs {
		ret = append(ret, string(ev.Value))
	}
	return
}

func (e *EtcdClient) SetAddr(key, value string) {
	_, err := e.Put(context.Background(), "laneIM/"+key, value)
	if err != nil {
		log.Fatalln("failed to set etcd:", err.Error())
		return
	}
}

func (e *EtcdClient) GetUserServerAddr(userid int64) string {
	rt, err := e.Get(context.Background(), "laneIM/us/"+strconv.FormatInt(userid, 36), clientv3.WithPrefix())
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return ""
	}
	for _, ev := range rt.Kvs {
		return string(ev.Value)
	}
	return ""
}

func (e *EtcdClient) SetUserServerAddr(userid int64, addr string) {
	_, err := e.Put(context.Background(), "laneIM/us/"+strconv.FormatInt(userid, 36), addr)
	if err != nil {
		log.Fatalln("failed to set etcd:", err.Error())
		return
	}
}

func (e *EtcdClient) GetUserRoomid(userid int64) UserRoom {
	rt, err := e.Get(context.Background(), e.UserIdToEtcdRoomKey(userid), clientv3.WithPrefix())
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
		return UserRoom{}
	}
	if len(rt.Kvs) == 0 {
		return UserRoom{}
	}
	userRoom := &UserRoom{}
	json.Unmarshal(rt.Kvs[0].Value, userRoom)
	return *userRoom
}

func (e *EtcdClient) AddUserRoomid(userid int64, roomid int64) bool {
	key := e.UserIdToEtcdRoomKey(userid)

	oldRoom := e.GetUserRoomid(userid)
	for _, v := range oldRoom.Id {
		if v == roomid {
			log.Printf("user[%d]  already in  room[%d]  \n", userid, roomid)
			return false
		}
	}
	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return false
	}

	oldRoom.Id = append(oldRoom.Id, roomid)
	newValue, err := json.Marshal(oldRoom)

	if err != nil {
		log.Println("json marshel err")
		return false
	}
	err = e.AtomicUpdate(key, string(oldValue), string(newValue))
	if err != nil {
		log.Println("failed:", err)
		return false
	}
	return true
}

func (e *EtcdClient) SetUserRoomid(userid int64, roomid []int64) bool {
	key := e.UserIdToEtcdRoomKey(userid)

	oldRoom := e.GetUserRoomid(userid)

	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return false
	}
	newRoom := UserRoom{
		Id: roomid,
	}
	newValue, err := json.Marshal(newRoom)
	if err != nil {
		log.Println("json marshel err")
		return false
	}

	err = e.AtomicUpdate(key, string(oldValue), string(newValue))
	if err != nil {
		log.Println("failed:", err)
		return false
	}
	return true
}

func (e *EtcdClient) DeleteUserRoom(userid int64, roomid int64) (rt bool) {
	key := e.UserIdToEtcdRoomKey(userid)

	oldRoom := e.GetUserRoomid(userid)
	oldValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return false
	}
	for i, v := range oldRoom.Id {
		if v == roomid {
			oldRoom.Id = append(oldRoom.Id[:i], oldRoom.Id[i+1:]...)
			rt = true
		}
	}

	if !rt {
		log.Printf("user[%d] do not have this room[%d]\n", userid, roomid)
		return false
	}

	newValue, err := json.Marshal(oldRoom)
	if err != nil {
		log.Println("json marshel err")
		return false
	}
	err = e.AtomicUpdate(key, string(oldValue), string(newValue))
	if err != nil {
		log.Println("failed")
		return false
	}
	return true
}

func (e *EtcdClient) UserIdToEtcdRoomKey(userid int64) string {
	return "laneIM/ur/" + strconv.FormatInt(userid, 36)
}

func (e *EtcdClient) UserIdToEtcdServerKey(userid int64) string {
	return "laneIM/us/" + strconv.FormatInt(userid, 36)
}

type UserRoom struct {
	Id []int64
}

func (e *EtcdClient) AtomicUpdate(key, old, new string) error {
	rt, err := conc.NewSTM(e.Client, func(s conc.STM) error {
		tmp := s.Get(key)
		// log.Printf("key[%v], old[%v], new[%v], tmp[%v]\n", key, old, new, tmp)
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
