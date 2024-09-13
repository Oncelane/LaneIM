package pkg_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"testing"
)

func TestRedis(t *testing.T) {
	e := pkg.NewEtcd(config.Etcd{Addr: []string{"127.0.0.1:2379"}})
	laneLog.Logger.Infoln("get: ", e.GetAddr("grpc:logic"))
}
