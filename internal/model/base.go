package model

import "time"

type BaseModel struct {
	ID        int       `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
}
