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
	conf = config.Logic{Type: "logic"}
	config.Init(&conf)
	l := logic.NewLogic(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	l.Close()

}
