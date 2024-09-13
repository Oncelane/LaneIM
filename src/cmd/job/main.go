package main

import (
	"flag"
	"laneIM/src/config"
	"laneIM/src/job"
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
	conf := config.Job{}
	config.Init(*ConfigPath, &conf)

	laneLog.InitLogger("job"+conf.Name, true)

	laneLog.Logger.Infof("job server start")
	j := job.NewJob(conf)
	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞等待信号

	j.Close()

}
