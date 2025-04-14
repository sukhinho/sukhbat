package models

import "time"

type Message struct {
	ID         uint   `gorm:"primaryKey"`
	Content    string `gorm:"type:text"`
	SenderID   uint
	ReceiverID uint
	CreatedAt  time.Time
}
