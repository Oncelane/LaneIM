package pkg

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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
		log.Fatalf("faile to connet to etcd: %v", err)
	}
	return &EtcdClient{etcdClient}
}

func (e *EtcdClient) GetAddr(key string) (ret []string) {
	rt, err := e.Get(context.Background(), key, clientv3.WithPrefix())
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
