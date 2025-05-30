package model

import (
	"time"
)

// KeywordGroup（正規化された技術キーワード）
type KeywordGroup struct {
	KeywordGroupID uint      `gorm:"primaryKey"`                                                        // キーワードグループID
	Name           string    `gorm:"unique;size:255;not null"`                                          // キーワード名（正規化）
	Type           string    `gorm:"type:enum('language','framework','skill','tool','other');not null"` // 分類
	CreatedAt      time.Time // 作成日時
	UpdatedAt      time.Time // 更新日時

	Words []KeyWord `gorm:"foreignKey:KeywordGroupID;references:KeywordGroupID"` // 表記ゆれ一覧
}
