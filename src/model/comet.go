package model

type CometMgr struct {
	CometAddr string    `gorm:"primary_key"`
	Rooms     []RoomMgr `gorm:"many2many:room_comets;"`
}
