package pkg_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
)

func TestRedis(t *testing.T) {
	e := pkg.NewEtcd(config.Etcd{Addr: []string{
		"127.0.0.1:51240",
		"127.0.0.1:51241",
		"127.0.0.1:51242"}})
	laneLog.Logger.Infoln("get: ", e.GetAddr("logic"))
	e.SetAddr("redis:1", "127.0.0.1:7001")
	e.SetAddr("redis:2", "127.0.0.1:7002")
	e.SetAddr("redis:3", "127.0.0.1:7003")
}
