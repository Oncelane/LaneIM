package canal

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/dao/rds"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/withlin/canal-go/client"
	"github.com/withlin/canal-go/protocol"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
)

type Canal struct {
	connector     *client.SimpleCanalConnector
	msgCh         chan *protocol.Message
	redis         *pkg.RedisClient
	db            *sql.SqlDB
	kafkaProducer *pkg.KafkaProducer
	conf          config.Canal
}

func NewCanal(conf config.Canal) *Canal {
	connector := client.NewSimpleCanalConnector(conf.CanalAddress, conf.CanalPort,
		conf.CanalName, conf.CanalPassword, conf.CanalDestination,
		conf.SoTimeOut, conf.IdleTimeOut)

	err := connector.Connect()
	if err != nil {
		laneLog.Logger.Infoln(err)
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
	c.db = sql.NewDB(conf.Mysql)
	laneLog.Logger.Infoln("start canal")
	return c

}

func (c *Canal) RunCanal() {
	for {

		message, err := c.connector.Get(100, nil, nil)
		if err != nil {
			laneLog.Logger.Infoln(err)
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
			laneLog.Logger.Infof("================> binlog[%s : %d],name[%s,%s], eventType: %s\n", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType())
			if len(header.GetTableName()) == 0 {
				continue
			}
			rediskey := strings.ReplaceAll(header.GetTableName(), "_", ":")[:len(header.GetTableName())-1]
			laneLog.Logger.Infof("redis key: %s\n", rediskey)
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
			c.HandleRowData(rowChange, eventType, rediskey)

		}
	}
}

func (c *Canal) HandleRowData(rowChange *pbe.RowChange, event pbe.EventType, redisKey string) {
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
			laneLog.Logger.Infoln("忽略", event)
		}
	}
}

func (c *Canal) HandleInsert(col []*pbe.Column, redisKey string) {
	rediskey := strings.Split(redisKey, ":")
	var rediskey0 = rediskey[0]
	var rediskey1 = rediskey[1]
	switch rediskey0 {
	case "room":
		switch rediskey1 {
		case "mgr":
			roomid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				laneLog.Logger.Infoln("fomat err")
				return
			}
			go c.db.MergeWriter.Do(redisKey, func() (any, error) {
				// sql
				room, err := c.db.RoomMgr(roomid)
				if err != nil {
					laneLog.Logger.Infoln("faild to sql roommgr")
				}
				// updata setex
				err = rds.SetEXRoomMgr(c.redis.Client, room)
				if err != nil {
					laneLog.Logger.Infoln("faild to sync ", redisKey, err)
				}

				laneLog.Logger.Infoln("sync", redisKey)
				return nil, nil
			})

		case "user":
			roomid, err := strconv.ParseInt(col[1].GetValue(), 10, 64)
			if err != nil {
				laneLog.Logger.Infoln("fomat err")
				return
			}
			go c.db.MergeWriter.Do(redisKey+col[1].GetValue(), func() (any, error) {

				// sql
				userids, err := c.db.RoomUserid(roomid)
				if err != nil {
					laneLog.Logger.Infoln("faild to sql ", redisKey)
				}

				// updata setex
				err = rds.SetEXRoomMgrUsers(c.redis.Client, roomid, userids)
				if err != nil {
					laneLog.Logger.Infoln("faild to sync ", redisKey, err)
				}
				laneLog.Logger.Infoln("sync", redisKey)
				return nil, nil
			})
		case "comet":
			roomid, err := strconv.ParseInt(col[1].GetValue(), 10, 64)
			if err != nil {
				laneLog.Logger.Infoln("fomat err")
				return
			}
			go c.db.MergeWriter.Do(redisKey+col[1].GetValue(), func() (any, error) {

				// sql
				comets, err := c.db.RoomComet(roomid)
				if err != nil {
					laneLog.Logger.Infoln("faild to sql roomid")
				}

				// updata setex
				err = rds.SetEXRoomCometBatch(c.redis.Client, roomid, comets)
				if err != nil {
					laneLog.Logger.Infoln("faild to sync ", redisKey, err)
				}
				laneLog.Logger.Infoln("sync", redisKey)
				return nil, nil
			})
		}
	case "user":
		switch rediskey1 {
		case "mgr":
			userid, err := strconv.ParseInt(col[0].GetValue(), 10, 64)
			if err != nil {
				laneLog.Logger.Infoln("fomat err")
				return
			}
			go c.db.MergeWriter.Do(redisKey, func() (any, error) {

				// sql
				user, err := c.db.UserMgr(userid)
				if err != nil {
					laneLog.Logger.Infoln("faild to sql user")
				}

				// updata setex
				err = rds.SetEXUserMgr(c.redis.Client, user)
				if err != nil {
					laneLog.Logger.Infoln("faild to sync ", redisKey, err)
				}

				laneLog.Logger.Infoln("sync", redisKey)
				return nil, nil

			})
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
