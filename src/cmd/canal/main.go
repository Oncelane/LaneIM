package main

import (
	"flag"
	"laneIM/src/canal"
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog"
	"time"
)

var (
	ConfigPath = flag.String("c", "config.yml", "path fo config.yml folder")
)

func main() {

	flag.Parse()
	conf := config.Canal{}
	config.Init(*ConfigPath, &conf)

	laneLog.InitLogger("canal"+conf.Name, true)

	laneLog.Logger.Infoln("[server] time", time.Duration(conf.Mysql.BatchWriter.MaxTime)*time.Millisecond, "count", conf.Mysql.BatchWriter.MaxCount)
	canal := canal.NewCanal(conf)
	go canal.RunCanal()
	go canal.RunReceive()
	select {}
}
