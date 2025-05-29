package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateEmailWorkTypeGroup はメールと業務グループの関連のサンプルデータを投入する。
func CreateEmailWorkTypeGroup(tx *gorm.DB) error {
	var err error

	emailWorkTypeGroups := []model.EmailWorkTypeGroup{
		// email001 (Java案件) の関連
		{
			EmailID:         "email001",
			WorkTypeGroupID: 2, // 基本設計
		},
		{
			EmailID:         "email001",
			WorkTypeGroupID: 3, // 詳細設計
		},
		{
			EmailID:         "email001",
			WorkTypeGroupID: 5, // バックエンド開発
		},
		{
			EmailID:         "email001",
			WorkTypeGroupID: 10, // 単体テスト
		},
		// email002 (React案件) の関連
		{
			EmailID:         "email002",
			WorkTypeGroupID: 3, // 詳細設計
		},
		{
			EmailID:         "email002",
			WorkTypeGroupID: 4, // フロントエンド開発
		},
		{
			EmailID:         "email002",
			WorkTypeGroupID: 19, // コードレビュー
		},
		// email003 (Python機械学習案件) の関連
		{
			EmailID:         "email003",
			WorkTypeGroupID: 1, // 要件定義
		},
		{
			EmailID:         "email003",
			WorkTypeGroupID: 2, // 基本設計
		},
		{
			EmailID:         "email003",
			WorkTypeGroupID: 5, // バックエンド開発
		},
		{
			EmailID:         "email003",
			WorkTypeGroupID: 7, // データベース設計
		},
		// email004 (Go案件) の関連
		{
			EmailID:         "email004",
			WorkTypeGroupID: 2, // 基本設計
		},
		{
			EmailID:         "email004",
			WorkTypeGroupID: 5, // バックエンド開発
		},
		{
			EmailID:         "email004",
			WorkTypeGroupID: 6, // API開発
		},
		{
			EmailID:         "email004",
			WorkTypeGroupID: 8, // インフラ構築
		},
		// email005 (フルスタック案件) の関連
		{
			EmailID:         "email005",
			WorkTypeGroupID: 3, // 詳細設計
		},
		{
			EmailID:         "email005",
			WorkTypeGroupID: 4, // フロントエンド開発
		},
		{
			EmailID:         "email005",
			WorkTypeGroupID: 5, // バックエンド開発
		},
		{
			EmailID:         "email005",
			WorkTypeGroupID: 6, // API開発
		},
		{
			EmailID:         "email005",
			WorkTypeGroupID: 8, // インフラ構築
		},
	}

	for _, emailWorkTypeGroup := range emailWorkTypeGroups {
		err := tx.Create(&emailWorkTypeGroup).Error
		if err != nil {
			return err
		}
	}

	return err
}
