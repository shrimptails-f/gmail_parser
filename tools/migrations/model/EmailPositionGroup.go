package model

// EmailPositionGroup（メールとポジショングループの中間）
type EmailPositionGroup struct {
	EmailID         string `gorm:"primaryKey;size:32"` // メールID
	PositionGroupID uint   `gorm:"primaryKey"`         // ポジショングループID
}
