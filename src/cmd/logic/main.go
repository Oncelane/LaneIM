package main

import (
	"laneIM/src/config"
	"laneIM/src/logic"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var conf config.Logic

func main() {
	conf = config.Logic{}
	config.Init("logic", &conf)
	log.Printf("logic server start by env:%+v", conf)
	l := logic.NewLogic(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	l.Close()

}
