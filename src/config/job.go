package config

import (
	"laneIM/src/pkg/util"
	"log"
)

type Job struct {
	Addr             string
	Name             string
	KafkaComsumer    KafkaComsumer
	Etcd             Etcd
	Redis            Redis
	BucketSize       int
	CometRoutineSize int
	Mysql            Mysql
}

func (c *Job) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	mysqlC := Mysql{}
	mysqlC.Default()
	*c = Job{
		Addr:             ip + ":50070",
		Name:             "0",
		KafkaComsumer:    DefaultKafkaComsumer(),
		Etcd:             DefaultEtcd(),
		Redis:            DefaultRedis(),
		BucketSize:       32,
		CometRoutineSize: 32,
		Mysql:            mysqlC,
	}
}
