package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/withlin/canal-go/client"
	"github.com/withlin/canal-go/protocol"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
)

type Canal struct {
	connector *client.SimpleCanalConnector
	msgCh     chan *protocol.Message
}

func NewCanal() *Canal {
	connector := client.NewSimpleCanalConnector("127.0.0.1", 11111,
		"canal", "canal", "laneIM",
		60000, 60*60*1000)
	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return &Canal{
		connector: connector,
		msgCh:     make(chan *protocol.Message, 128),
	}
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
		log.Println("todo: write to redis")
		printEntry(message.Entries)
	}
}

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
	canal := NewCanal()
	go canal.RunCanal()
	go canal.RunReceive()
	select {}
}

func printEntry(entrys []pbe.Entry) {

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
