package cql_test

import (
	"fmt"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao/cql"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg/laneLog"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants"
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

// 制造SendMsgBatchReq测试数据
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
	start := time.Now()
	wait := sync.WaitGroup{}
	wait.Add(t.N)
	for range t.N {
		go func(t *testing.B) {
			defer wait.Done()
			count, err = db.PageChatMessageUnReadCountByMessageid(1005, -1, offlineTimpStamp)
		}(t)
	}
	wait.Wait()
	if err != nil {
		t.Fatal(err)
	}
	laneLog.Logger.Debugln("online, query count", t.N, "unread message count:", count, "spand time:", time.Since(start))
}
func BenchmarkGetPageChatMessageUnReadCountByMessageidWith_Pool(t *testing.B) {
	// laneLog.Logger.Debugln("time.Now=", time.Now())
	// laneLog.Logger.Debugln("timestamp.Now=", timestamppb.Now().AsTime())

	db := GetDB_B(t)
	GeneratorMessage_B(t, db, 100, 1005)
	// 60秒前下线
	offlineTimpStamp := time.Now().Add(time.Second * -60)
	var count int
	var err error
	start := time.Now()
	wait := sync.WaitGroup{}
	wait.Add(t.N)

	p, _ := ants.NewPool(150)

	for range t.N {
		err = p.Submit(func() {
			defer wait.Done()
			count, err = db.PageChatMessageUnReadCountByMessageid(1005, -1, offlineTimpStamp)
		})
		if err != nil {
			wait.Done()
		}
	}
	wait.Wait()
	if err != nil {
		t.Fatal(err)
	}
	laneLog.Logger.Debugln("online, query count", t.N, "unread message count:", count, "spand time:", time.Since(start))
}

// ChatMessage 数据模型
type ChatMessageMysql struct {
	GroupID   int64     `gorm:"primaryKey;column:group_id"`
	UserID    int64     `gorm:"primaryKey;column:user_id"`
	UserSeq   int64     `gorm:"column:user_seq"`
	Timestamp time.Time `gorm:"primaryKey;column:timestamp"`
	MessageID int64     `gorm:"primaryKey;column:message_id"`
	Content   []byte    `gorm:"column:content"`
}

// ConnectDB 用于连接数据库并自动迁移表
func ConnectDB_Test(t *testing.B) *sql.SqlDB {
	t.Helper()
	db := sql.NewDB(config.DefaultMysql())
	// 自动迁移数据库表
	if err := db.DB.AutoMigrate(&ChatMessageMysql{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	return db
}

func doGeneratorMessage_Mysql(t *testing.B, db *sql.SqlDB, messageSize int, roomid int64) error {
	t.Helper()
	messages := make([]ChatMessageMysql, messageSize)
	for i := range messages {
		messages[i] = ChatMessageMysql{
			MessageID: int64(i),
			UserID:    int64(i),
			GroupID:   roomid,
			UserSeq:   0,
			Timestamp: time.Now(),
			Content:   []byte(fmt.Sprintf("i am message from %d", i)),
		}
	}
	if err := db.DB.Create(&messages).Error; err != nil {
		return err
	}
	return nil
}

func PageChatMessageUnReadCountByMessageid_Mysql(db *sql.SqlDB, groupID int64, lastMessageID int64, lastMsgUnix time.Time) (count int, err error) {
	// 创建查询语句
	// debugSql := `
	// 	SELECT COUNT(*)
	// 	FROM chat_message_mysqls
	// 	WHERE group_id = ? AND timestamp >= ?;`

	// laneLog.Logger.Debugln("sql = ", debugSql)

	// 执行查询并获取结果
	err = db.DB.Model(&ChatMessageMysql{}).Select("count(*)").Where("group_id = ? AND timestamp >= ?", groupID, lastMsgUnix).Scan(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

func BenchmarkGetPageChatMessageUnReadCountByMessageidWithMysql(t *testing.B) {
	// laneLog.Logger.Debugln("time.Now=", time.Now())
	// laneLog.Logger.Debugln("timestamp.Now=", timestamppb.Now().AsTime())

	db := ConnectDB_Test(t)
	mysqlDB, _ := db.DB.DB()
	defer mysqlDB.Close()
	doGeneratorMessage_Mysql(t, db, 100, 1005)
	var count int
	var err error
	offlineTimpStamp := time.Now().Add(time.Second * -60)
	start := time.Now()
	// 假设 t.N 是你需要的并发任务数量
	tN := 150 // 示例数量

	// 创建一个协程池，限制最大并发数为 tN
	p, err := ants.NewPool(tN)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer p.Release()

	// 等待组用于等待所有任务完成
	wait := sync.WaitGroup{}
	wait.Add(t.N)

	for i := 0; i < t.N; i++ {
		// 将任务提交到协程池
		err := p.Submit(func() {
			defer wait.Done()
			count, err = PageChatMessageUnReadCountByMessageid_Mysql(db, 1005, -1, offlineTimpStamp)
			if err != nil {
				log.Printf("Error fetching unread message count: %v", err)
			}
		})
		if err != nil {
			log.Printf("Failed to submit task: %v", err)
			wait.Done() // 确保在错误情况下也调用 Done
		}
	}

	// 等待所有任务完成
	wait.Wait()
	if err != nil {
		t.Fatal(err)
	}
	laneLog.Logger.Debugln("online, query count", t.N, "unread message count:", count, "spand time:", time.Since(start))
}
