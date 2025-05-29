package model

import (
	"time"
)

// Position（職種マスタ）
type Position struct {
	ID        uint      `gorm:"primaryKey"`              // ポジションID
	Name      string    `gorm:"unique;size:50;not null"` // ポジション名（例: PL）
	CreatedAt time.Time // 作成日時
	UpdatedAt time.Time // 更新日時
}
