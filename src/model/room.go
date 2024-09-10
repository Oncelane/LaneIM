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
		&UserComet{},
		// &UserRoom{},
		&UserOnline{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

// RoomMgr 模型
type RoomMgr struct {
	RoomID int64     `gorm:"primary_key"`
	Users  []UserMgr `gorm:"many2many:room_users;"`
}

// RoomOnline 模型
type RoomOnline struct {
	RoomID      int64 `gorm:"primary_key"`
	OnlineCount int   `gorm:"not null"`
}

// RoomComet 模型
type RoomComet struct {
	RoomID    int64  `gorm:"primary_key"`
	CometAddr string `gorm:"type:varchar(255);"`
}

// RoomUserid 模型
type RoomUser struct {
	RoomID int64 `gorm:"primary_key"`
	UserID int64 `gorm:"primary_key"`
}
