package main

import (
	"laneIM/src/config"
	"laneIM/src/logic"
	"os"
	"os/signal"
	"syscall"
)

var conf config.Logic

func main() {
	conf = config.Logic{
		Addr: "127.0.0.1:50050",
		Name: "logic1",
		KafkaProducer: config.KafkaProducer{
			Addr: []string{"127.0.0.1:9092"},
		},
		Etcd: config.Etcd{
			Addr: []string{"127.0.0.1:2379"},
		},
	}
	l := logic.NewLogic(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	l.Close()

}
