package main

import (
	"laneIM/src/comet"
	"laneIM/src/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.Comet{Type: "comet"}
	config.Init(&conf)
	log.Printf("[%s] server start by env:%+v", conf.GetType()+conf.GetName(), conf)
	c := comet.NewSerivceComet(conf)

	// 启动websocket服务
	http.HandleFunc("/ws", c.ServeHTTP)
	go http.ListenAndServe(":40051", nil)

	// 等待信号
	log.Println("wait ctrl+c")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号
	log.Println("close comet")
	c.Close()

}
