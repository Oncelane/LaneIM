package model

// UserMgr 模型
type UserMgr struct {
	UserID int64 `gorm:"primary_key;"`
}

// UserComet 模型
type UserComet struct {
	UserID    int64  `gorm:"primary_key"`
	CometAddr string `gorm:"type:varchar(255);not null"`
}

// UserRoom 模型
type UserRoom struct {
	UserID int64 `gorm:"primary_key"`
	RoomID int64
}

// UserOnline 模型
type UserOnline struct {
	UserID int64 `gorm:"primary_key"`
	Online bool  `gorm:"not null"`
}
