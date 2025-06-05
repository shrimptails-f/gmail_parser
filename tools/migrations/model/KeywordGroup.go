package model

import (
	"time"
)

// KeywordGroup（正規化された技術キーワード）
type KeywordGroup struct {
	KeywordGroupID uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"size:255;not null;uniqueIndex:idx_name_type"`
	Type           string `gorm:"type:enum('language','framework','must','want','other');not null;uniqueIndex:idx_name_type"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// 循環してて完全に積んでるのでコメントアウト
	// WordLinks []KeywordGroupWordLink `gorm:"foreignKey:KeywordGroupID;references:KeywordGroupID"`
}
