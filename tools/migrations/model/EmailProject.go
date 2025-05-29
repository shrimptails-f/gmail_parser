package model

import (
	"time"
)

// EmailProject（案件メール専用情報）
type EmailProject struct {
	EmailID         string    `gorm:"primaryKey;size:32"` // メールID（emails.idと同じ）
	ProjectTitle    *string   `gorm:"size:255"`           // 案件名
	EntryTimings    *string   `gorm:"type:text"`          // 入場時期（カンマ or JSON）
	EndTiming       *string   `gorm:"size:255"`           // 終了時期
	WorkLocation    *string   `gorm:"size:255;index"`     // 勤務場所
	PriceFrom       *int      `gorm:"type:int"`           // 単価FROM
	PriceTo         *int      `gorm:"type:int"`           // 単価TO
	RemoteType      *string   `gorm:"size:50"`            // リモート区分
	RemoteFrequency *string   `gorm:"size:255"`           // リモート頻度
	CreatedAt       time.Time // 作成日時
	UpdatedAt       time.Time // 更新日時

	// リレーション
	Email Email `gorm:"foreignKey:EmailID;references:ID"` // 親メール
}
