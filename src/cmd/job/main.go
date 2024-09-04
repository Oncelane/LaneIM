package main

import (
	"laneIM/src/config"
	"laneIM/src/job"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.Job{Type: "job"}
	config.Init(&conf)
	j := job.NewJob(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	j.Close()

}
