package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateKeywordGroup はキーワードグループのサンプルデータを投入する。
func CreateKeywordGroup(tx *gorm.DB) error {
	var err error
	keywordGroups := []model.KeywordGroup{
		{
			ID:   1,
			Name: "Java",
			Type: "language",
		},
		{
			ID:   2,
			Name: "Python",
			Type: "language",
		},
		{
			ID:   3,
			Name: "JavaScript",
			Type: "language",
		},
		{
			ID:   4,
			Name: "TypeScript",
			Type: "language",
		},
		{
			ID:   5,
			Name: "Go",
			Type: "language",
		},
		{
			ID:   6,
			Name: "Spring Boot",
			Type: "framework",
		},
		{
			ID:   7,
			Name: "React",
			Type: "framework",
		},
		{
			ID:   8,
			Name: "Vue.js",
			Type: "framework",
		},
		{
			ID:   9,
			Name: "Django",
			Type: "framework",
		},
		{
			ID:   10,
			Name: "AWS",
			Type: "tool",
		},
		{
			ID:   11,
			Name: "Docker",
			Type: "tool",
		},
		{
			ID:   12,
			Name: "Kubernetes",
			Type: "tool",
		},
		{
			ID:   13,
			Name: "MySQL",
			Type: "tool",
		},
		{
			ID:   14,
			Name: "PostgreSQL",
			Type: "tool",
		},
		{
			ID:   15,
			Name: "設計",
			Type: "skill",
		},
		{
			ID:   16,
			Name: "テスト",
			Type: "skill",
		},
		{
			ID:   17,
			Name: "レビュー",
			Type: "skill",
		},
	}

	for _, keywordGroup := range keywordGroups {
		err := tx.Create(&keywordGroup).Error
		if err != nil {
			return err
		}
	}

	return err
}
