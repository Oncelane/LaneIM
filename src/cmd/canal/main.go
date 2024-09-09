package main

import (
	"flag"
	"fmt"
	"laneIM/src/config"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/withlin/canal-go/client"
	"github.com/withlin/canal-go/protocol"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Canal struct {
	connector     *client.SimpleCanalConnector
	msgCh         chan *protocol.Message
	redis         *pkg.RedisClient
	db            *gorm.DB
	kafkaProducer *pkg.KafkaProducer
	conf          config.Canal
}

func NewCanal(conf config.Canal) *Canal {
	connector := client.NewSimpleCanalConnector(conf.CanalAddress, conf.CanalPort,
		conf.CanalName, conf.CanalPassword, conf.CanalDestination,
		conf.SoTimeOut, conf.IdleTimeOut)

	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// connector.RollBack(-1)
	connector.Subscribe(conf.Subscribe)
	c := &Canal{
		connector:     connector,
		msgCh:         make(chan *protocol.Message, conf.MsgChSize),
		redis:         pkg.NewRedisClient(conf.Redis),
		kafkaProducer: pkg.NewKafkaProducer(conf.KafkaProducer),
		conf:          conf,
	}
	c.db = sql.DB(conf.Mysql)
	log.Println("start canal")
	return c

}

func (c *Canal) RunCanal() {
	for {

		message, err := c.connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			// fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "===没有数据了===")
			continue
		}

		c.msgCh <- message

	}
}

func (c *Canal) RunReceive() {
	for message := range c.msgCh {
		c.printEntry(message.Entries)
	}
}

var (
	ConfigPath = flag.String("c", "config.yml", "path fo config.yml folder")
)

func main() {

	// 192.168.199.17 替换成你的canal server的地址
	// example 替换成-e canal.destinations=example 你自己定义的名字
	//  该字段名字在 canal\conf\example\meta.dat 文件中，NewSimpleCanalConnector函数参数配置，也在文件中
	/**
	  NewSimpleCanalConnector 参数说明
	    client.NewSimpleCanalConnector("Canal服务端地址", "Canal服务端端口", "Canal服务端用户名", "Canal服务端密码", "Canal服务端destination", 60000, 60*60*1000)
	    Canal服务端地址：canal服务搭建地址IP
	    Canal服务端端口：canal\conf\canal.properties文件中
	    Canal服务端用户名、密码：canal\conf\example\instance.properties 文件中
	    Canal服务端destination ：canal\conf\example\meta.dat 文件中
	*/

	// https://github.com/alibaba/canal/wiki/AdminGuide
	//mysql 数据解析关注的表，Perl正则表达式.
	//
	//多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
	//
	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//  2.  canal schema下所有表： canal\\..*
	//  3.  canal下的以canal打头的表：canal\\.canal.*
	//  4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)
	flag.Parse()
	conf := config.Canal{}
	config.Init(*ConfigPath, &conf)
	canal := NewCanal(conf)
	go canal.RunCanal()
	go canal.RunReceive()
	select {}
}

func (c *Canal) printEntry(entrys []pbe.Entry) {
	for i := range entrys {
		if entrys[i].GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entrys[i].GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		err := proto.Unmarshal(entrys[i].GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entrys[i].GetHeader()
			fmt.Printf("================> binlog[%s : %d],name[%s,%s], eventType: %s\n", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType())
			if len(header.GetTableName()) == 0 {
				continue
			}
			rediskey := strings.ReplaceAll(header.GetTableName(), "_", ":")[:len(header.GetTableName())-1]
			fmt.Printf("redis key: %s\n", rediskey)
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_DELETE {
					printColumn(rowData.GetBeforeColumns())
				} else if eventType == pbe.EventType_INSERT {
					printColumn(rowData.GetAfterColumns())
				} else {
					fmt.Println("-------> before")
					printColumn(rowData.GetBeforeColumns())
					fmt.Println("-------> after")
					printColumn(rowData.GetAfterColumns())
				}
			}
			c.HandleRowData(rowChange, eventType, strings.Split(rediskey, ":"))

		}
	}
}

func (c *Canal) HandleRowData(rowChange *pbe.RowChange, event pbe.EventType, redisKey []string) {
	for _, col := range rowChange.GetRowDatas() {
		switch event {
		// case pbe.EventType_CREATE:
		case pbe.EventType_UPDATE:
			c.HandleInsert(col.GetBeforeColumns(), redisKey)
		case pbe.EventType_INSERT:
			c.HandleInsert(col.GetAfterColumns(), redisKey)
		case pbe.EventType_DELETE:
			c.HandleInsert(col.GetBeforeColumns(), redisKey)
		default:
			log.Println("忽略", event)
		}
	}
}

func (c *Canal) HandleInsert(col []*pbe.Column, redisKey []string) {
	switch redisKey[0] {
	case "room":
		switch redisKey[1] {
		case "mgr":

			// check exist
			exist, err := rds.EXAllRoomid(c.redis.Client)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			if !exist {
				return
			}

			// sql
			roomids, err := sql.AllRoomid(c.db)
			if err != nil {
				log.Println("faild to sql roomid")
			}

			// updata setex
			err = rds.SetEXAllRoomid(c.redis.Client, roomids)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			log.Println("sync room:mgr")

		case "online":
			roomid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}

			// check exist
			exist, err := rds.EXRoomOnline(c.redis.Client, roomid)
			if err != nil {
				log.Println("faild to check ", err)
				return
			}
			if !exist {
				return
			}

			//sql
			count, err := strconv.ParseInt(col[1].GetValue(), 10, 64)
			if err != nil {
				log.Println("faild to decode onlinecount", col[1].GetValue(), err)
			}
			err = rds.SetEXRoomOnlie(c.redis.Client, roomid, int(count))
			if err != nil {
				log.Println("faild to sync room:online", err)
			}
			log.Println("success", roomid)
		case "userid":
			roomid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}
			// check exist
			exist, err := rds.EXRoomUser(c.redis.Client, roomid)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			if !exist {
				return
			}

			// sql
			userids, err := sql.RoomUserid(c.db, roomid)
			if err != nil {
				log.Println("faild to sql roomid")
			}

			// updata setex
			err = rds.SetEXRoomUser(c.redis.Client, roomid, userids)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			log.Println("sync room:user", roomid)
		case "comet":
			roomid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}
			// check exist
			exist, err := rds.EXRoomComet(c.redis.Client, roomid)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			if !exist {
				return
			}

			// sql
			comets, err := sql.RoomComet(c.db, roomid)
			if err != nil {
				log.Println("faild to sql roomid")
			}

			// updata setex
			err = rds.SetEXRoomComet(c.redis.Client, roomid, comets)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			log.Println("sync room:comet", roomid)
		}
	case "user":
		switch redisKey[1] {
		case "mgr":

			// check exist
			exist, err := rds.EXAllUserid(c.redis.Client)
			if err != nil {
				log.Println("faild to sync user:mgr", err)
			}
			if !exist {
				return
			}

			// sql
			userids, err := sql.AllUserid(c.db)
			if err != nil {
				log.Println("faild to sql user")
			}

			// updata setex
			err = rds.SetEXAllUserid(c.redis.Client, userids)
			if err != nil {
				log.Println("faild to sync user:mgr", err)
			}
			log.Println("sync user:mgr")

		case "online":
			roomid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}

			// check exist
			exist, err := rds.EXRoomOnline(c.redis.Client, roomid)
			if err != nil {
				log.Println("faild to check ", err)
				return
			}
			if !exist {
				return
			}

			//sql
			count, err := strconv.ParseInt(col[1].GetValue(), 10, 64)
			if err != nil {
				log.Println("faild to decode onlinecount", col[1].GetValue(), err)
			}
			err = rds.SetEXRoomOnlie(c.redis.Client, roomid, int(count))
			if err != nil {
				log.Println("faild to sync room:online", err)
			}
			log.Println("success", roomid)
		case "room":
			userid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}
			// check exist
			exist, err := rds.EXUserRoom(c.redis.Client, userid)
			if err != nil {
				log.Println("faild to sync user:room", err)
			}
			if !exist {
				return
			}

			// sql
			roomids, err := sql.UserRoom(c.db, userid)
			if err != nil {
				log.Println("faild to sql user:room")
			}

			// updata setex
			err = rds.SetEXUserRoom(c.redis.Client, userid, roomids)
			if err != nil {
				log.Println("faild to sync room:mgr", err)
			}
			log.Println("sync user:room", userid)
		case "comet":
			userid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				log.Println("fomat err")
				return
			}
			// check exist
			exist, err := rds.EXUserComet(c.redis.Client, userid)
			if err != nil {
				log.Println("faild to sync user:comet", err)
			}
			if !exist {
				return
			}

			// sql
			cometAddr, err := sql.UserComet(c.db, userid)
			if err != nil {
				log.Println("faild to sql user")
			}

			// updata setex
			err = rds.SetEXUserComet(c.redis.Client, userid, cometAddr)
			if err != nil {
				log.Println("faild to sync user:comet", err)
			}
			log.Println("sync user:comet", userid)
		}
	}
}

func printColumn(columns []*pbe.Column) {
	for _, col := range columns {
		fmt.Printf("%s : %s  update= %t\n", col.GetName(), col.GetValue(), col.GetUpdated())
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
