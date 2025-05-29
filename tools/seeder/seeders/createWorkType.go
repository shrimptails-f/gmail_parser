package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateWorkType は業務タイプのサンプルデータを投入する。
func CreateWorkType(tx *gorm.DB) error {
	var err error
	workTypes := []model.WorkType{
		{
			ID:   1,
			Name: "要件定義",
		},
		{
			ID:   2,
			Name: "基本設計",
		},
		{
			ID:   3,
			Name: "詳細設計",
		},
		{
			ID:   4,
			Name: "フロントエンド実装",
		},
		{
			ID:   5,
			Name: "バックエンド実装",
		},
		{
			ID:   6,
			Name: "API実装",
		},
		{
			ID:   7,
			Name: "データベース設計",
		},
		{
			ID:   8,
			Name: "インフラ構築",
		},
		{
			ID:   9,
			Name: "単体テスト",
		},
		{
			ID:   10,
			Name: "結合テスト",
		},
		{
			ID:   11,
			Name: "システムテスト",
		},
		{
			ID:   12,
			Name: "運用保守",
		},
		{
			ID:   13,
			Name: "パフォーマンスチューニング",
		},
		{
			ID:   14,
			Name: "セキュリティ対応",
		},
		{
			ID:   15,
			Name: "ドキュメント作成",
		},
		{
			ID:   16,
			Name: "コードレビュー",
		},
		{
			ID:   17,
			Name: "技術調査",
		},
		{
			ID:   18,
			Name: "プロジェクト管理",
		},
		{
			ID:   19,
			Name: "チームマネジメント",
		},
		{
			ID:   20,
			Name: "顧客折衝",
		},
	}

	for _, workType := range workTypes {
		err := tx.Create(&workType).Error
		if err != nil {
			return err
		}
	}

	return err
}
