package main

import (
	"flag"
	"laneIM/src/config"
	"laneIM/src/logic"
	"laneIM/src/pkg/laneLog"
	"os"
	"os/signal"
	"syscall"
)

var (
	ConfigPath = flag.String("c", "config.yml", "path fo config.yml folder")
)

func main() {
	flag.Parse()
	conf := config.Logic{}
	config.Init(*ConfigPath, &conf)

	laneLog.InitLogger("logic"+conf.Name, true)

	laneLog.Logger.Infoln("[server] logic server start on ", conf.Addr)
	l := logic.NewLogic(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	l.Close()

}
