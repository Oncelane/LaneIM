package model

import (
	"log"

	"gorm.io/gorm"
)

func Init(db *gorm.DB) {

	// 自动迁移
	err := db.AutoMigrate(
		&RoomMgr{},
		&UserMgr{},
		&RoomOnline{},
		&RoomComet{},
		&RoomUserid{},
		&UserComet{},
		&UserRoom{},
		&UserOnline{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

// RoomMgr 模型
type RoomMgr struct {
	RoomID int64 `gorm:"primaryKey;"`
}

// RoomOnline 模型
type RoomOnline struct {
	RoomID      int64 `gorm:"primaryKey"`
	OnlineCount int   `gorm:"not null"`
}

// RoomComet 模型
type RoomComet struct {
	ID        uint   `gorm:"primaryKey;"`
	RoomID    int64  `gorm:"not null"`
	CometAddr string `gorm:"type:varchar(255);not null;unique"`
}

// RoomUserid 模型
type RoomUserid struct {
	ID     uint  `gorm:"primaryKey;"`
	RoomID int64 `gorm:"not null"`
	UserID int64 `gorm:"unique;not null"`
}
