package model

import (
	"time"
)

// Email（メール基本情報）
type Email struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"` // オートインクリメントID
	GmailID      string    `gorm:"size:32;index"`            // GメールID
	Subject      string    `gorm:"type:text;not null"`       // 件名
	SenderName   string    `gorm:"size:255"`                 // 差出人名
	SenderEmail  string    `gorm:"size:255;index"`           // メールアドレス
	ReceivedDate time.Time `gorm:"index"`                    // 受信日
	Body         *string   `gorm:"type:longtext"`            // 本文
	Category     string    `gorm:"size:50;index"`            // 種別（案件 / 人材提案）

	CreatedAt time.Time // 作成日時
	UpdatedAt time.Time // 更新日時

	// 子テーブル
	EmailProject        *EmailProject        `gorm:"foreignKey:EmailID;references:ID"` // 案件情報（1対1）
	EmailCandidate      *EmailCandidate      `gorm:"foreignKey:EmailID;references:ID"` // 人材情報（1対1）
	EmailKeywordGroups  []EmailKeywordGroup  `gorm:"foreignKey:EmailID;references:ID"` // 技術キーワード（1対多）
	EmailPositionGroups []EmailPositionGroup `gorm:"foreignKey:EmailID;references:ID"` // ポジション（1対多）
	EmailWorkTypeGroups []EmailWorkTypeGroup `gorm:"foreignKey:EmailID;references:ID"` // 業務内容（1対多）
}
