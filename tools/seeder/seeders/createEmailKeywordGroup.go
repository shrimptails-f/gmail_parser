package seeders

import (
	"business/tools/migrations/model"
	"time"

	"gorm.io/gorm"
)

// CreateEmailKeywordGroup はメールとキーワードグループの関連のサンプルデータを投入する。
func CreateEmailKeywordGroup(tx *gorm.DB) error {
	var err error

	emailKeywordGroups := []model.EmailKeywordGroup{
		// email001 (Java案件) の関連
		{
			EmailID:        1,
			KeywordGroupID: 1, // Java
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        1,
			KeywordGroupID: 6, // Spring Boot
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        1,
			KeywordGroupID: 13, // MySQL
			Type:           "WANT",
			CreatedAt:      time.Now(),
		},
		// email002 (React案件) の関連
		{
			EmailID:        2,
			KeywordGroupID: 7, // React
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        2,
			KeywordGroupID: 4, // TypeScript
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        2,
			KeywordGroupID: 3, // JavaScript
			Type:           "LANGUAGE",
			CreatedAt:      time.Now(),
		},
		// email003 (Python機械学習案件) の関連
		{
			EmailID:        3,
			KeywordGroupID: 2, // Python
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        3,
			KeywordGroupID: 9, // Django
			Type:           "FRAMEWORK",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        3,
			KeywordGroupID: 10, // AWS
			Type:           "WANT",
			CreatedAt:      time.Now(),
		},
		// email004 (Go案件) の関連
		{
			EmailID:        4,
			KeywordGroupID: 5, // Go
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        4,
			KeywordGroupID: 11, // Docker
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        4,
			KeywordGroupID: 12, // Kubernetes
			Type:           "WANT",
			CreatedAt:      time.Now(),
		},
		// email005 (フルスタック案件) の関連
		{
			EmailID:        5,
			KeywordGroupID: 7, // React
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        5,
			KeywordGroupID: 3, // JavaScript
			Type:           "LANGUAGE",
			CreatedAt:      time.Now(),
		},
		{
			EmailID:        5,
			KeywordGroupID: 10, // AWS
			Type:           "MUST",
			CreatedAt:      time.Now(),
		},
	}

	for _, emailKeywordGroup := range emailKeywordGroups {
		err := tx.Create(&emailKeywordGroup).Error
		if err != nil {
			return err
		}
	}

	return err
}
