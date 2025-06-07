package model

import (
	"time"
)

// EmailKeywordGroup（メールとキーワードの多対多）
type EmailKeywordGroup struct {
	EmailID        uint `gorm:"not null;uniqueIndex:idx_email_keyword_type"`
	KeywordGroupID uint `gorm:"not null;uniqueIndex:idx_email_keyword_type"`
	CreatedAt      time.Time

	// 循環してて完全に積んでるのでコメントアウト
	// KeywordGroup KeywordGroup `gorm:"foreignKey:KeywordGroupID;references:KeywordGroupID"`
}
