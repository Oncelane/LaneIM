package main

import (
	"laneIM/src/config"
	"laneIM/src/job"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.Job{
		Addr: "127.0.0.50052",
		Name: "job1",
		Etcd: config.Etcd{
			Addr: []string{"127.0.0.1:2379"},
		},
		KafkaComsumer: config.KafkaComsumer{
			Addr: []string{"127.0.0.1:9092"},
		},
	}
	j := job.NewJob(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	j.Close()

}
