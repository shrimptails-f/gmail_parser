package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateEmailWorkType はメールと業務タイプの関連のサンプルデータを投入する。
func CreateEmailWorkType(tx *gorm.DB) error {
	var err error

	emailWorkTypes := []model.EmailWorkType{
		// email001 (Java案件) の関連
		{
			EmailID:    "email001",
			WorkTypeID: 2, // 基本設計
		},
		{
			EmailID:    "email001",
			WorkTypeID: 3, // 詳細設計
		},
		{
			EmailID:    "email001",
			WorkTypeID: 5, // バックエンド実装
		},
		{
			EmailID:    "email001",
			WorkTypeID: 9, // 単体テスト
		},
		// email002 (React案件) の関連
		{
			EmailID:    "email002",
			WorkTypeID: 3, // 詳細設計
		},
		{
			EmailID:    "email002",
			WorkTypeID: 4, // フロントエンド実装
		},
		{
			EmailID:    "email002",
			WorkTypeID: 16, // コードレビュー
		},
		// email003 (Python機械学習案件) の関連
		{
			EmailID:    "email003",
			WorkTypeID: 1, // 要件定義
		},
		{
			EmailID:    "email003",
			WorkTypeID: 2, // 基本設計
		},
		{
			EmailID:    "email003",
			WorkTypeID: 5, // バックエンド実装
		},
		{
			EmailID:    "email003",
			WorkTypeID: 17, // 技術調査
		},
		// email004 (Go案件) の関連
		{
			EmailID:    "email004",
			WorkTypeID: 2, // 基本設計
		},
		{
			EmailID:    "email004",
			WorkTypeID: 5, // バックエンド実装
		},
		{
			EmailID:    "email004",
			WorkTypeID: 6, // API実装
		},
		{
			EmailID:    "email004",
			WorkTypeID: 8, // インフラ構築
		},
		// email005 (フルスタック案件) の関連
		{
			EmailID:    "email005",
			WorkTypeID: 3, // 詳細設計
		},
		{
			EmailID:    "email005",
			WorkTypeID: 4, // フロントエンド実装
		},
		{
			EmailID:    "email005",
			WorkTypeID: 5, // バックエンド実装
		},
		{
			EmailID:    "email005",
			WorkTypeID: 6, // API実装
		},
		{
			EmailID:    "email005",
			WorkTypeID: 8, // インフラ構築
		},
	}

	for _, emailWorkType := range emailWorkTypes {
		err := tx.Create(&emailWorkType).Error
		if err != nil {
			return err
		}
	}

	return err
}
