package pkg

import (
	"errors"
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog"
	"log"

	ec "github.com/Oncelane/laneEtcd/src/client"
	"github.com/Oncelane/laneEtcd/src/kvraft"
	"github.com/Oncelane/laneEtcd/src/pkg/laneConfig"
)

type EtcdClient struct {
	etcd *ec.Clerk
}

var (
	ErrOpNotAtomic = errors.New("data change not atomic op")
	ErrFaild       = errors.New("op faild for some reason")
)

func NewEtcd(conf config.Etcd) *EtcdClient {
	etcdClient := ec.MakeClerk(laneConfig.Clerk{
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
	strs := make([]string, len(rt))
	for i := range rt {
		strs[i] = string(rt[i])
	}
	return strs
}

func (e *EtcdClient) SetAddr(key, addr string) (cancel func()) {
	// err := e.etcd.Put("laneIM:"+key, []byte(addr), 0)
	cancel = e.etcd.WatchDog("laneIM:"+key, []byte(addr))
	laneLog.Logger.Infof("[laneEtcd] registe %s:%s\n", key, addr)
	return
}

func (e *EtcdClient) SetAddr_WithoutWatch(key, addr string) error {
	err := e.etcd.Put("laneIM:"+key, []byte(addr), 0)
	// cancel = e.etcd.WatchDog("laneIM:"+key, []byte(addr))
	laneLog.Logger.Infof("[laneEtcd] registe %s:%s\n", key, addr)
	return err
}

func (e *EtcdClient) DelAddr(key, addr string) {
	err := e.etcd.Delete("laneIM:" + key)
	if err != nil {
		log.Fatalln("failed to get etcd:", err.Error())
	}
	laneLog.Logger.Infoln("delete addr:", addr)
}
