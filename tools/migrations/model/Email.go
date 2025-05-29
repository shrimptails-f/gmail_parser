package model

import (
	"time"
)

// Email（メール基本情報）
type Email struct {
	ID           string    `gorm:"primaryKey;size:32"` // メールID
	Subject      string    `gorm:"type:text;not null"` // 件名
	SenderName   string    `gorm:"size:255"`           // 差出人名
	SenderEmail  string    `gorm:"size:255;index"`     // メールアドレス
	ReceivedDate time.Time `gorm:"index"`              // 受信日
	Body         *string   `gorm:"type:longtext"`      // 本文
	Category     string    `gorm:"size:50;index"`      // 種別（案件 / 人材提案）

	CreatedAt time.Time // 作成日時
	UpdatedAt time.Time // 更新日時

	// リレーション
	KeywordGroups []KeywordGroup `gorm:"many2many:email_keyword_groups;"` // 技術キーワード（多対多）
	Positions     []Position     `gorm:"many2many:email_positions;"`      // ポジション（多対多）
	WorkTypes     []WorkType     `gorm:"many2many:email_work_types;"`     // 業務内容（多対多）

	// 子テーブル
	EmailProject   *EmailProject   `gorm:"foreignKey:EmailID;references:ID"` // 案件情報（1対1）
	EmailCandidate *EmailCandidate `gorm:"foreignKey:EmailID;references:ID"` // 人材情報（1対1）
}
