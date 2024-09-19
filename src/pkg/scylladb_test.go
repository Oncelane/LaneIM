package pkg

import (
	"fmt"
	"laneIM/src/pkg/laneLog"
	"strconv"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v3"
)

var keyspace string = "examples"

// metadata specifies table name and columns it must be in sync with schema.
var chatMessageMetadata = table.Metadata{
	Name:    "examples.messages",
	Columns: []string{"group_id", "user_id", "timestamp", "message_id", "content"},
	PartKey: []string{"group_id"},
	SortKey: []string{"timestamp", "message_id"},
}

// personTable allows for simple CRUD operations based on personMetadata.
var chatMessageTable = table.New(chatMessageMetadata)

// Person represents a row in person table.
// Field names are converted to snake case by default, no need to add special tags.
// A field will not be persisted by adding the `db:"-"` tag or making it unexported.
type ChatMessage struct {
	GroupID   int64
	UserID    int64
	Timestamp time.Time
	MessageID int64
	Content   []byte
}

func TestDB(t *testing.T) {
	// Create gocql cluster.
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		t.Fatal(err)
	}

	session.ExecStmt(`DROP KEYSPACE examples`)
	err = session.ExecStmt(fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		keyspace,
	))
	if err != nil {
		t.Fatal("create keyspace:", err)
	}

	err = session.ExecStmt(`CREATE TABLE examples.messages (
    group_id bigint,
    message_id bigint,
    timestamp TIMESTAMP,
    user_id bigint,
    content BLOB,
    PRIMARY KEY (group_id, message_id)
) WITH CLUSTERING ORDER BY ( message_id DESC);`)
	if err != nil {
		t.Fatal("create table:", err)
	}
	groupID := 1005
	var lastMessageID int64 // 上一页最后一条消息的消息id

	// { // 单插

	// 	// Insert 100 records.
	// 	for i := 0; i < 100000; i++ {
	// 		lastTimestamp = time.Now()
	// 		p := ChatMessage{
	// 			GroupID:   int64(groupID),
	// 			UserID:    int64(i),
	// 			Timestamp: lastTimestamp,
	// 			MessageID: int64(i),
	// 			Content:   []byte(strconv.FormatInt(int64(i), 10)),
	// 		}
	// 		q := session.Query(chatMessageTable.Insert()).BindStruct(p).Consistency(gocql.LocalOne)
	// 		if err := q.ExecRelease(); err != nil {
	// 			laneLog.Logger.Fatal(err)
	// 		}
	// 	}
	// }

	{ // 多插入
		// 创建 Batch
		batch := session.NewBatch(gocql.LoggedBatch)
		// 创建 Batch
		// batch := gocql.NewBatch(gocql.LoggedBatch)
		batch.Cons = gocql.LocalOne

		// Insert 100 records.
		for i := 0; i < 100; i++ {
			lastMessageID = int64(i)
			data := ChatMessage{
				GroupID:   int64(groupID),
				UserID:    int64(i),
				Timestamp: time.Now(),
				MessageID: int64(i),
				Content:   []byte(strconv.FormatInt(int64(i), 10)),
			}

			insertChatQryStmt, _ := qb.Insert("examples.messages").Columns(chatMessageTable.Metadata().Columns...).ToCql()
			batch.Query(insertChatQryStmt,
				PchatChatMessageToSlice(data)...)

		}
		lastMessageID++

		if err := session.ExecuteBatch(batch); err != nil {
			laneLog.Logger.Fatalln(err)
		}
	}

	{ // 分页查询示例
		var limit int = 10 // 每页限制
		curpage := 0
		for range 11 {

			debugSql := fmt.Sprintf(`
		SELECT *
		FROM examples.messages
		WHERE group_id = %s
		AND message_id < %d
		LIMIT %s`, strconv.FormatInt(int64(groupID), 10), lastMessageID, strconv.FormatInt(int64(limit), 10))
			// laneLog.Logger.Debugln(debugSql)
			// 假设你想获取第一页
			iter := session.Query(debugSql, nil).Iter()

			var oneMessage ChatMessage
			// Scan(&oneMessage.GroupID, &oneMessage.Timestamp, &oneMessage.MessageID, &oneMessage.Content, &oneMessage.UserID)
			for iter.StructScan(&oneMessage) {
				fmt.Printf("Message: %s, User: %d, Timestamp: %s\n", string(oneMessage.Content), oneMessage.UserID, oneMessage.Timestamp.String())
				lastMessageID = oneMessage.UserID // 更新最后一条消息的时间戳
			}

			if err := iter.Close(); err != nil {
				laneLog.Logger.Fatalf("Error during iteration: %v", err)
			}
			laneLog.Logger.Infof("==============page %d===========", curpage)
			curpage++
		}
	}

	// {
	// 	session.Query().PageSize()
	// }
}

func PchatChatMessageToSlice(in ChatMessage) []any {
	return []any{
		in.GroupID,
		in.UserID,
		in.Timestamp,
		in.MessageID,
		in.Content,
	}
}

func mustParseUUID(s string) gocql.UUID {
	u, err := gocql.ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}
