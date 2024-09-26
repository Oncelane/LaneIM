package cql_test

import (
	"fmt"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao/cql"
	"testing"
	"time"
)

func TestAddAndPage(t *testing.T) {
	conf := config.DefaultScyllaDB()
	conf.Keyspace = "examples"
	db := cql.NewCqlDB(conf)
	db.ResetKeySpace()
	db.InitChatMessage()

	num := 10000
	limit := 100
	testMessage(t, db, num, 1005)
	testMessage(t, db, num, 404)
	testPageMessage(t, db, limit, 1005)
	testPageMessage(t, db, limit, 404)
}

func testMessage(t *testing.T, db *cql.ScyllaDB, num int, room int64) {
	msgs := make([]*msg.SendMsgReq, num)
	for i := range num {
		// data := ChatMessage{
		// 	GroupID:   int64(groupID),
		// 	UserID:    int64(i),
		// 	Timestamp: time.Now(),
		// 	MessageID: int64(i),
		// 	Content:   []byte(strconv.FormatInt(int64(i), 10)),
		// }
		msgs[i] = &msg.SendMsgReq{
			Messageid: int64(i),
			Userid:    int64(i),
			Roomid:    room,
			Userseq:   0,
			Timeunix:  time.Now().Unix(),
			Data:      []byte(fmt.Sprintf("i am message from %d", i)),
		}
		// laneLog.Logger.Debugf("check data:%s", string(msgs[i].Data))
	}
	err := db.AddChatMessageBatch(&msg.SendMsgBatchReq{
		Msgs: msgs,
	})
	if err != nil {
		t.Fatal(err)
	}

}

func testPageMessage(t *testing.T, db *cql.ScyllaDB, limit int, room int64) {
	lastMsgid, lastMsgUnix, err := db.QueryLatestGroupMessageid(room)
	if err != nil {
		t.Fatal(err)
	}

	curpage := 0
	for {
		// laneLog.Logger.Debugln("last room message id =", lastMsgid)
		rt, validSize, err := db.PageChatMessageByMessageid(int64(1005), lastMsgid, lastMsgUnix, limit)
		// select {}
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < validSize; i++ {
			// laneLog.Logger.Infof("msg:%+s", rt[i].Content)
		}
		if validSize == 0 {
			break
		}
		curpage++
		lastMsgid = rt[validSize-1].MessageID
		// laneLog.Logger.Infof("room %d==============page %d===========", room, curpage)

	}

}
