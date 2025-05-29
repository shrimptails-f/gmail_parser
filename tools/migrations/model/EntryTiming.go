package model

import (
	"time"
)

type EntryTiming struct {
	ID             uint   `gorm:"primaryKey"`       // ID
	EmailProjectID string `gorm:"size:32;index"`    // 紐づく案件メールID
	StartDate      string `gorm:"size:20;not null"` // 入場日（例: "2025/06/01"）
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
