package config

import (
	"laneIM/src/pkg/util"
	"log"
)

type Logic struct {
	Id            int
	UbuntuIP      string
	WindowIP      string
	GrpcPort      string
	Name          string
	KafkaProducer KafkaProducer
	Etcd          Etcd
	Mysql         Mysql
	ScyllaDB      ScyllaDB
}

func (c *Logic) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	*c = Logic{
		Id:       0,
		UbuntuIP: "172.xx",
		WindowIP: "192.xx",
		GrpcPort: ":50060",
		Name:     "0",
		KafkaProducer: KafkaProducer{
			Addr: []string{ip + ":9092"},
		},
		Etcd:     DefaultEtcd(),
		Mysql:    DefaultMysql(),
		ScyllaDB: DefaultScyllaDB(),
	}
}
