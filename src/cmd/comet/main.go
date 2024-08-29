package main

import (
	"laneIM/proto/logic"
	"laneIM/src/comet"
	"laneIM/src/config"
	"log"
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

	log.Println("send test message")
	testSend(c)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号
	log.Println("close comet")
	c.Close()

}

func testSend(comet *comet.Comet) {
	msg := &logic.SendMsgReq{
		Data:   []byte("hello laneIM!"),
		Roomid: 1005,
		Path:   "brodcast/test",
		Userid: 21,
	}

	comet.LogictRoom(msg)
}
