package dao

import (
	"laneIM/src/dao/cql"
	"laneIM/src/model"
	"log"

	"gorm.io/gorm"
)

func Init(db *gorm.DB, cql *cql.ScyllaDB) {
	if db != nil {
		// 自动迁移
		err := db.AutoMigrate(
			&model.RoomMgr{},
			&model.UserMgr{},
			&model.CometMgr{},

			// &RoomMessage{},
		)
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	}

	if cql != nil {
		cql.InitChatMessage()
	}

}
