package model

import "time"

type RoomMessage struct {
	MessageId  int64     `gorm:"primary_key"`
	RoomId     int64     `gorm:"not null"`
	SenderId   int64     `gorm:"not null"`
	SenderDate time.Time `gorm:"not null"`
	Content    string    `gorm:"type:text;not null"`
}
