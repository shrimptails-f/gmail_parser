package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateWorkTypeGroup は業務グループのサンプルデータを投入する。
func CreateWorkTypeGroup(tx *gorm.DB) error {
	var err error

	workTypeGroups := []model.WorkTypeGroup{
		{ID: 1, Name: "要件定義"},
		{ID: 2, Name: "基本設計"},
		{ID: 3, Name: "詳細設計"},
		{ID: 4, Name: "フロントエンド開発"},
		{ID: 5, Name: "バックエンド開発"},
		{ID: 6, Name: "API開発"},
		{ID: 7, Name: "データベース設計"},
		{ID: 8, Name: "インフラ構築"},
		{ID: 9, Name: "テスト設計"},
		{ID: 10, Name: "単体テスト"},
		{ID: 11, Name: "結合テスト"},
		{ID: 12, Name: "システムテスト"},
		{ID: 13, Name: "運用保守"},
		{ID: 14, Name: "パフォーマンス改善"},
		{ID: 15, Name: "セキュリティ対応"},
		{ID: 16, Name: "CI/CD構築"},
		{ID: 17, Name: "監視設定"},
		{ID: 18, Name: "ドキュメント作成"},
		{ID: 19, Name: "コードレビュー"},
		{ID: 20, Name: "リファクタリング"},
	}

	for _, workTypeGroup := range workTypeGroups {
		err := tx.Create(&workTypeGroup).Error
		if err != nil {
			return err
		}
	}

	return err
}
