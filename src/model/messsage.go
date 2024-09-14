package model

import "time"

type GroupMessage struct {
	GroupId     int64     `gorm:"primary_key;not null"`
	SenderId    int64     `gorm:"not null"`
	MessageId   int64     `gorm:"not null"`
	MessageDate time.Time `gorm:"not null"`
	Content     string    `gorm:"type:text;not null"`
}
