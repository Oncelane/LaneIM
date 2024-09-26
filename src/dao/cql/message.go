package cql

import (
	"fmt"
	"laneIM/proto/msg"
	"laneIM/src/model"
	"laneIM/src/pkg/laneLog"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
)

func (s *ScyllaDB) InitChatMessage() {
	err := s.DB.ExecStmt(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.messages  (
	group_id bigint,
	message_id bigint,
	timestamp TIMESTAMP,
	user_id bigint,
	user_seq bigint,
	content BLOB,
	PRIMARY KEY (group_id, timestamp, message_id)
) WITH CLUSTERING ORDER BY ( timestamp DESC, message_id DESC);`, s.conf.Keyspace))
	if err != nil {
		laneLog.Logger.Fatalln("create table:", err)
	}

	// metadata specifies table name and columns it must be in sync with schema.
	var chatMessageMetadata = table.Metadata{
		Name:    fmt.Sprintf(`%s.messages`, s.conf.Keyspace),
		Columns: []string{"group_id", "user_id", "user_seq", "timestamp", "message_id", "content"},
		PartKey: []string{"group_id"},
		SortKey: []string{"timestamp,message_id"},
	}

	// personTable allows for simple CRUD operations based on personMetadata.
	s.chatMessageTable = table.New(chatMessageMetadata)
}

func (s *ScyllaDB) AddChatMessageBatch(in *msg.SendMsgBatchReq) error {
	// start := time.Now()
	{ // 多插入
		// 创建 Batch
		batch := s.DB.NewBatch(gocql.LoggedBatch)
		// 创建 Batch
		// batch := gocql.NewBatch(gocql.LoggedBatch)
		batch.Cons = gocql.LocalOne
		insertChatQryStmt, _ := qb.Insert(fmt.Sprintf("%s.messages", s.conf.Keyspace)).Columns(s.chatMessageTable.Metadata().Columns...).ToCql()

		for _, m := range in.Msgs {
			// data := model.ChatMessage{
			// 	GroupID:   m.Roomid,
			// 	UserID:    m.Userid,
			// 	UserSeq:   m.Userseq,
			// 	Timestamp: time.Unix(m.Timeunix, 0),
			// 	MessageID: m.Messageid,
			// 	Content:   m.Data,
			// }
			// laneLog.Logger.Debugf("check in.Msgs data:%s", string(m.Data))
			batch.Query(insertChatQryStmt,
				m.Roomid, m.Userid, m.Userseq, time.Unix(m.Timeunix, 0), m.Messageid, m.Data)
		}

		if err := s.DB.ExecuteBatch(batch); err != nil {
			laneLog.Logger.Fatalln(err)
			return err
		}
	}
	// laneLog.Logger.Debugln("batch insert spand time:", time.Since(start))
	return nil
}

func (s *ScyllaDB) PageChatMessageByMessageid(groupID int64, lastMessageID int64, lastMsgUnix int64, limit int) ([]model.ChatMessage, int, error) {
	// 分页查询示例
	debugSql := fmt.Sprintf(`
	SELECT *
	FROM %s.messages
	WHERE group_id = %d
	AND timestamp <= '%s'
	LIMIT %d`, s.conf.Keyspace, groupID, time.Unix(lastMsgUnix, 0).Format("2006-01-02T15:04:05.000Z"), limit)
	laneLog.Logger.Debugln("sql = ", debugSql)
	iter := s.DB.Query(debugSql, nil).Iter()
	rt := make([]model.ChatMessage, limit+1)
	// Scan(&oneMessage.GroupID, &oneMessage.Timestamp, &oneMessage.MessageID, &oneMessage.Content, &oneMessage.UserID)
	var index = 0
	for iter.StructScan(&rt[index]) {
		index++
		// laneLog.Logger.Debugf("Message: %s, User: %d, Timestamp: %s\n", string(oneMessage.Content), oneMessage.UserID, oneMessage.Timestamp.String())
	}

	if err := iter.Close(); err != nil {
		laneLog.Logger.Fatalf("Error during iteration: %v", err)
		return nil, 0, err
	}
	if index == 0 {
		return nil, 0, nil
	}

	return rt, index, nil
}

func (s *ScyllaDB) QueryLatestGroupMessageid(groupID int64) (lastMessageID, lastMsgUnix int64, err error) {
	// 分页查询示例
	debugSql := fmt.Sprintf(`
	SELECT *
	FROM %s.messages
	WHERE group_id = %d
	LIMIT 1`, s.conf.Keyspace, groupID)
	// laneLog.Logger.Debugln("sql = ", debugSql)
	var oneMessage []model.ChatMessage
	// Scan(&oneMessage.GroupID, &oneMessage.Timestamp, &oneMessage.MessageID, &oneMessage.Content, &oneMessage.UserID)
	err = s.DB.Query(debugSql, nil).Select(&oneMessage)
	if err != nil {
		return
	}
	if len(oneMessage) != 0 {
		return oneMessage[0].MessageID, oneMessage[0].Timestamp.Unix(), nil
	}

	return -1, -1, gocql.ErrNotFound
}
