package model

// UserMgr 模型
type UserMgr struct {
	UserID    int64     `gorm:"primary_key"`
	Rooms     []RoomMgr `gorm:"many2many:room_users;"`
	CometAddr string    `gorm:"type:varchar(255);not null"`
	Online    bool      `gorm:"not null"`
}
