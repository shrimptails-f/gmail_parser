package model

// EmailPosition（メールとポジションの中間）
type EmailPosition struct {
	EmailID    string `gorm:"primaryKey;size:32"` // メールID
	PositionID uint   `gorm:"primaryKey"`         // ポジションID
}
