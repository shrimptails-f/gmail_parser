package model

// EmailWorkTypeGroup（メールと業務グループの中間）
type EmailWorkTypeGroup struct {
	EmailID         string `gorm:"primaryKey;size:32"` // メールID
	WorkTypeGroupID uint   `gorm:"primaryKey"`         // 業務グループID
}
