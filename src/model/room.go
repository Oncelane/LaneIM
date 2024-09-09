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
	RoomID    int64  `gorm:"primaryKey"`
	CometAddr string `gorm:"primaryKey;type:varchar(255);not null"`
}

// RoomUserid 模型
type RoomUserid struct {
	RoomID int64 `gorm:"primaryKey"`
	UserID int64 `gorm:"primaryKey"`
}
