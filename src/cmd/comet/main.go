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
	conf := config.Comet{}
	config.Init("comet", &conf)
	log.Printf("comet server start by env:%+v", conf)
	c := comet.NewSerivceComet(conf)

	// 启动websocket服务
	http.HandleFunc("/ws", c.ServeHTTP)
	log.Println("listening websocket", conf.WebsocketAddr)
	go http.ListenAndServe(conf.WebsocketAddr, nil)

	// 等待信号
	log.Println("wait ctrl+c")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号
	log.Println("close comet")
	c.Close()

}
