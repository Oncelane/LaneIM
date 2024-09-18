package pkg_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
)

var etcd pkg.EtcdClient

func init() {
	conf := config.Etcd{
		Addr: []string{"127.0.0.1:51240", "127.0.0.1:51241", "127.0.0.1:51242"},
	}
	etcd = *pkg.NewEtcd(conf)
}

func TestEtcd(t *testing.T) {
	etcd.SetAddr("test1", "1")

	etcd.SetAddr("test2", "2")

	etcd.SetAddr("test3", "3")
	laneLog.Logger.Infoln(etcd.GetAddr("test"))
}
