package pkg

import (
	"errors"
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog"
	"log"

	"github.com/Oncelane/laneEtcd/src/kvraft"
	"github.com/Oncelane/laneEtcd/src/pkg/laneConfig"
)

type EtcdClient struct {
	etcd *kvraft.Clerk
}

var (
	ErrOpNotAtomic = errors.New("data change not atomic op")
	ErrFaild       = errors.New("op faild for some reason")
)

func NewEtcd(conf config.Etcd) *EtcdClient {
	etcdClient := kvraft.MakeClerk(laneConfig.Clerk{
		EtcdAddrs: conf.Addr,
	})
	return &EtcdClient{etcdClient}
}

func (e *EtcdClient) GetAddr(key string) []string {
	rt, err := e.etcd.GetWithPrefix("laneIM:" + key)
	if err != nil && err != kvraft.ErrNil {
		log.Fatalln("failed to get etcd:", err.Error())
		return nil
	}
	return rt
}

func (e *EtcdClient) SetAddr(key, addr string) {
	err := e.etcd.Put("laneIM:"+key, addr)
	if err != nil {
		log.Fatalln("failed to set etcd:", err.Error())
		return
	}
	laneLog.Logger.Infof("registe %s:%s\n", key, addr)
}

func (e *EtcdClient) DelAddr(key, addr string) {
	err := e.etcd.Delete("laneIM:" + key)
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
	}
	laneLog.Logger.Infoln("delete addr:", addr)
}
