package model

// EmailWorkType（メールと業務の中間）
type EmailWorkType struct {
	EmailID    string `gorm:"primaryKey;size:32"` // メールID
	WorkTypeID uint   `gorm:"primaryKey"`         // 業務ID
}
