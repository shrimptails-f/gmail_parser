package model

import (
	"time"
)

// EmailKeywordGroup（メールとキーワードの多対多）
type EmailKeywordGroup struct {
	EmailID        uint `gorm:"not null;"`
	KeywordGroupID uint `gorm:"not null;"`
	CreatedAt      time.Time

	// 循環してて完全に積んでるのでコメントアウト
	// KeywordGroup KeywordGroup `gorm:"foreignKey:KeywordGroupID;references:KeywordGroupID"`
}
