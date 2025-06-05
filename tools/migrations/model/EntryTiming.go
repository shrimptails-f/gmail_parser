package model

import (
	"time"
)

type EntryTiming struct {
	EmailID   uint   // ID
	StartDate string `gorm:";size:20;not null"` // 入場日（例: "2025/06/01"）
	CreatedAt time.Time
	UpdatedAt time.Time
}
