package model

type CometMgr struct {
	CometAddr string    `gorm:"primary_key"`
	Valid     bool      `gorm:"not null"`
	Rooms     []RoomMgr `gorm:"many2many:room_comets;"`
}
