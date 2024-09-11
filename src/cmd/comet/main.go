package main

import (
	"flag"
	"laneIM/src/comet"
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog.go"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	ConfigPath = flag.String("c", "config.yml", "path fo config.yml folder")
)

func main() {
	flag.Parse()
	conf := config.Comet{}
	config.Init(*ConfigPath, &conf)

	laneLog.InitLogger("comet"+conf.Name, true)

	laneLog.Logger.Infof("comet server start")
	c := comet.NewSerivceComet(conf)

	// 启动websocket服务
	http.HandleFunc("/ws", c.ServeHTTP)
	laneLog.Logger.Infoln("listening websocket", conf.WebsocketAddr)
	go http.ListenAndServe(conf.WebsocketAddr, nil)

	// 等待信号
	laneLog.Logger.Infoln("wait ctrl+c")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号
	laneLog.Logger.Infoln("close comet")
	c.Close()

}
