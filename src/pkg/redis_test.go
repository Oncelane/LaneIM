package pkg_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	e := pkg.NewEtcd(config.Etcd{Addr: []string{"127.0.0.1:2379"}})
	log.Println("get: ", e.GetAddr("grpc:logic"))
}
