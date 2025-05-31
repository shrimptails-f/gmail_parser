package model

// EmailPositionGroup（メールとポジショングループの中間）
type EmailPositionGroup struct {
	EmailID         uint `gorm:"primaryKey"` // メールID
	PositionGroupID uint `gorm:"primaryKey"` // ポジショングループID
}
