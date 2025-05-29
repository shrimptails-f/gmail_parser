package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateEmailPosition はメールとポジションの関連のサンプルデータを投入する。
func CreateEmailPosition(tx *gorm.DB) error {
	var err error

	emailPositions := []model.EmailPosition{
		// email001 (Java案件) の関連
		{
			EmailID:    "email001",
			PositionID: 2, // SE
		},
		{
			EmailID:    "email001",
			PositionID: 3, // PL
		},
		// email002 (React案件) の関連
		{
			EmailID:    "email002",
			PositionID: 8, // フロントエンドエンジニア
		},
		{
			EmailID:    "email002",
			PositionID: 2, // SE
		},
		// email003 (Python機械学習案件) の関連
		{
			EmailID:    "email003",
			PositionID: 15, // 機械学習エンジニア
		},
		{
			EmailID:    "email003",
			PositionID: 14, // データエンジニア
		},
		// email004 (Go案件) の関連
		{
			EmailID:    "email004",
			PositionID: 9, // バックエンドエンジニア
		},
		{
			EmailID:    "email004",
			PositionID: 6, // アーキテクト
		},
		// email005 (フルスタック案件) の関連
		{
			EmailID:    "email005",
			PositionID: 10, // フルスタックエンジニア
		},
		{
			EmailID:    "email005",
			PositionID: 8, // フロントエンドエンジニア
		},
		{
			EmailID:    "email005",
			PositionID: 9, // バックエンドエンジニア
		},
	}

	for _, emailPosition := range emailPositions {
		err := tx.Create(&emailPosition).Error
		if err != nil {
			return err
		}
	}

	return err
}
