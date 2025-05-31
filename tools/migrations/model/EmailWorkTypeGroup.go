package model

// EmailWorkTypeGroup（メールと業務グループの中間）
type EmailWorkTypeGroup struct {
	EmailID         uint `gorm:"primaryKey"` // メールID
	WorkTypeGroupID uint `gorm:"primaryKey"` // 業務グループID
}
