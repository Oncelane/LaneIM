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
	conf := config.Comet{
		Addr: "127.0.0.1:50051",
		Name: "comet1",
		Etcd: config.Etcd{
			Addr: []string{"127.0.0.1:2379"},
		},
	}
	c := comet.NewSerivceComet(conf)

	// 启动websocket服务
	http.HandleFunc("/ws", c.ServeHTTP)
	http.ListenAndServe(":40051", nil)

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号
	log.Println("close comet")
	c.Close()

}
