package model

import (
	"time"
)

// EmailKeywordGroup（メールとキーワードの多対多）
type EmailKeywordGroup struct {
	EmailID        uint      `gorm:"primaryKey;size:32"`                                               // メールID
	KeywordGroupID uint      `gorm:"primaryKey"`                                                       // キーワードグループID
	Type           string    `gorm:"primaryKey;type:enum('MUST','WANT','LANGUAGE','FRAMEWORK');index"` // スキル種別
	CreatedAt      time.Time // 登録日時
}
