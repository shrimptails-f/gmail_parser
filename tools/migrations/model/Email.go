package model

import (
	"time"
)

// Email（メール基本情報）
type Email struct {
	ID              string    `gorm:"primaryKey;size:32"` // メールID
	Subject         string    `gorm:"type:text;not null"` // 件名
	SenderName      string    `gorm:"size:255"`           // 差出人名
	SenderEmail     string    `gorm:"size:255;index"`     // メールアドレス
	ReceivedDate    time.Time `gorm:"index"`              // 受信日
	Body            *string   `gorm:"type:text"`          // 本文
	Category        *string   `gorm:"size:50;index"`      // メール区分（案件 / 人材提案）
	ProjectTitle    *string   `gorm:"size:255"`           // 案件名
	EntryTimings    *string   `gorm:"type:text"`          // 入場時期・開始時期
	EndTiming       *string   `gorm:"size:255"`           // 終了時期
	WorkLocation    *string   `gorm:"size:255;index"`     // 勤務場所
	PriceFrom       *int      `gorm:"type:int"`           // 単価FROM
	PriceTo         *int      `gorm:"type:int"`           // 単価TO
	RemoteType      *string   `gorm:"size:50"`            // リモートワーク区分
	RemoteFrequency *string   `gorm:"size:255"`           // リモートワークの頻度
	CreatedAt       time.Time // 作成日時
	UpdatedAt       time.Time // 更新日時

	KeywordGroups []KeywordGroup `gorm:"many2many:email_keyword_groups;"` // 技術キーワード（多対多）
	Positions     []Position     `gorm:"many2many:email_positions;"`      // ポジション（多対多）
	WorkTypes     []WorkType     `gorm:"many2many:email_work_types;"`     // 業務内容（多対多）
}
