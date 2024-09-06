package model

// UserMgr 模型
type UserMgr struct {
	ID     uint   `gorm:"primary_key;auto_increment"`
	UserID uint64 `gorm:"unique;not null"`
}

// UserComet 模型
type UserComet struct {
	UserID    uint64 `gorm:"primaryKey"`
	CometInfo string `gorm:"type:varchar(255);not null"`
}

// UserRoom 模型
type UserRoom struct {
	UserID uint64 `gorm:"primaryKey"`
	RoomID uint64 `gorm:"primaryKey"`
}

// UserOnline 模型
type UserOnline struct {
	UserID uint64 `gorm:"primaryKey"`
	Online bool   `gorm:"not null"`
}
