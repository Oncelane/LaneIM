package cql_test

import (
	"fmt"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao/cql"
	"laneIM/src/pkg/laneLog"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddAndPage(t *testing.T) {
	db := GetDB_T(t)

	num := 10000
	limit := 100
	GeneratorMessage_T(t, db, num, 1005)
	GeneratorMessage_T(t, db, num, 404)
	testPageMessage(t, db, limit, 1005)
	testPageMessage(t, db, limit, 404)
}

func GetDB_T(t *testing.T) *cql.ScyllaDB {
	t.Helper()
	conf := config.DefaultScyllaDB()
	conf.Keyspace = "examples"
	db := cql.NewCqlDB(conf)
	db.ResetKeySpace()
	db.InitChatMessage()
	if db == nil {
		t.Fatal("faild to connect scyllaDB")
	}
	return db
}

func GetDB_B(t *testing.B) *cql.ScyllaDB {
	t.Helper()
	conf := config.DefaultScyllaDB()
	conf.Keyspace = "examples"
	db := cql.NewCqlDB(conf)
	db.ResetKeySpace()
	db.InitChatMessage()
	if db == nil {
		t.Fatal("faild to connect scyllaDB")
	}
	return db
}
func doGeneratorMessage(db *cql.ScyllaDB, num int, room int64) error {
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
			Timeunix:  timestamppb.Now(),
			Data:      []byte(fmt.Sprintf("i am message from %d", i)),
		}
		// laneLog.Logger.Debugf("check data:%s", string(msgs[i].Data))
	}
	err := db.AddChatMessageBatch(&msg.SendMsgBatchReq{
		Msgs: msgs,
	})
	return err
}
func GeneratorMessage_T(t *testing.T, db *cql.ScyllaDB, num int, room int64) {
	t.Helper()
	err := doGeneratorMessage(db, num, room)
	if err != nil {
		t.Fatal(err)
	}
}
func GeneratorMessage_B(t *testing.B, db *cql.ScyllaDB, num int, room int64) {
	t.Helper()
	err := doGeneratorMessage(db, num, room)
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

func BenchmarkGetPageChatMessageUnReadCountByMessageid(t *testing.B) {
	// laneLog.Logger.Debugln("time.Now=", time.Now())
	// laneLog.Logger.Debugln("timestamp.Now=", timestamppb.Now().AsTime())

	db := GetDB_B(t)
	GeneratorMessage_B(t, db, 100, 1005)
	// 60秒前下线
	offlineTimpStamp := time.Now().Add(time.Second * -60)
	var count int
	var err error
	for range t.N {
		count, err = db.PageChatMessageUnReadCountByMessageid(1005, -1, offlineTimpStamp)
		if err != nil {
			t.Fatal(err)
		}
	}
	laneLog.Logger.Debugln("online, query unread message count:", count)
}
