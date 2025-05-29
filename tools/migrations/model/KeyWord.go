package model

import (
	"time"
)

// KeyWord（表記ゆれ）
type KeyWord struct {
	ID             uint      `gorm:"primaryKey"`        // 表記ゆれID
	KeywordGroupID uint      `gorm:"not null"`          // 対応するキーワードグループID
	Word           string    `gorm:"size:255;not null"` // 表記ゆれ文字列
	CreatedAt      time.Time // 作成日時
	UpdatedAt      time.Time // 更新日時
}
