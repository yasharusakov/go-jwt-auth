package entity

import "time"

type RefreshToken struct {
	ID        int       `gorm:"type:serial;primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:integer;unique;not null;" json:"user_id"`
	Token     string    `gorm:"type:text;unique;not null;" json:"token"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoCreateTime" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
