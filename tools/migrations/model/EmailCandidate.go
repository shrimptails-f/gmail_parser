package model

import (
	"time"
)

// EmailCandidate（人材提案メール専用情報）
type EmailCandidate struct {
	EmailID          uint      `gorm:"primaryKey"` // メールID
	CandidateName    *string   `gorm:"size:255"`   // 人材名（仮）
	ExperienceYears  *int      `gorm:"type:int"`   // 経験年数
	SkillsSummary    *string   `gorm:"type:text"`  // 自己紹介・スキルまとめ
	AvailabilityDate *string   `gorm:"size:255"`   // 参画可能日
	CreatedAt        time.Time // 作成日時
	UpdatedAt        time.Time // 更新日時

	// リレーション
	Email Email `gorm:"foreignKey:EmailID;references:ID"` // 親メール
}
