package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreateKeyWord はキーワード（表記ゆれ）のサンプルデータを投入する。
func CreateKeyWord(tx *gorm.DB) error {
	var err error
	keyWords := []model.KeyWord{
		// Java関連の表記ゆれ
		{
			ID:             1,
			KeywordGroupID: 1,
			Word:           "Java",
		},
		{
			ID:             2,
			KeywordGroupID: 1,
			Word:           "JAVA",
		},
		{
			ID:             3,
			KeywordGroupID: 1,
			Word:           "java",
		},
		// Python関連の表記ゆれ
		{
			ID:             4,
			KeywordGroupID: 2,
			Word:           "Python",
		},
		{
			ID:             5,
			KeywordGroupID: 2,
			Word:           "python",
		},
		{
			ID:             6,
			KeywordGroupID: 2,
			Word:           "PYTHON",
		},
		// JavaScript関連の表記ゆれ
		{
			ID:             7,
			KeywordGroupID: 3,
			Word:           "JavaScript",
		},
		{
			ID:             8,
			KeywordGroupID: 3,
			Word:           "Javascript",
		},
		{
			ID:             9,
			KeywordGroupID: 3,
			Word:           "JS",
		},
		{
			ID:             10,
			KeywordGroupID: 3,
			Word:           "js",
		},
		// TypeScript関連の表記ゆれ
		{
			ID:             11,
			KeywordGroupID: 4,
			Word:           "TypeScript",
		},
		{
			ID:             12,
			KeywordGroupID: 4,
			Word:           "Typescript",
		},
		{
			ID:             13,
			KeywordGroupID: 4,
			Word:           "TS",
		},
		{
			ID:             14,
			KeywordGroupID: 4,
			Word:           "ts",
		},
		// Go関連の表記ゆれ
		{
			ID:             15,
			KeywordGroupID: 5,
			Word:           "Go",
		},
		{
			ID:             16,
			KeywordGroupID: 5,
			Word:           "Golang",
		},
		{
			ID:             17,
			KeywordGroupID: 5,
			Word:           "golang",
		},
		// Spring Boot関連の表記ゆれ
		{
			ID:             18,
			KeywordGroupID: 6,
			Word:           "Spring Boot",
		},
		{
			ID:             19,
			KeywordGroupID: 6,
			Word:           "SpringBoot",
		},
		{
			ID:             20,
			KeywordGroupID: 6,
			Word:           "Spring",
		},
		// React関連の表記ゆれ
		{
			ID:             21,
			KeywordGroupID: 7,
			Word:           "React",
		},
		{
			ID:             22,
			KeywordGroupID: 7,
			Word:           "React.js",
		},
		{
			ID:             23,
			KeywordGroupID: 7,
			Word:           "ReactJS",
		},
		// Vue.js関連の表記ゆれ
		{
			ID:             24,
			KeywordGroupID: 8,
			Word:           "Vue.js",
		},
		{
			ID:             25,
			KeywordGroupID: 8,
			Word:           "Vue",
		},
		{
			ID:             26,
			KeywordGroupID: 8,
			Word:           "VueJS",
		},
		// Django関連の表記ゆれ
		{
			ID:             27,
			KeywordGroupID: 9,
			Word:           "Django",
		},
		{
			ID:             28,
			KeywordGroupID: 9,
			Word:           "django",
		},
		// AWS関連の表記ゆれ
		{
			ID:             29,
			KeywordGroupID: 10,
			Word:           "AWS",
		},
		{
			ID:             30,
			KeywordGroupID: 10,
			Word:           "Amazon Web Services",
		},
		{
			ID:             31,
			KeywordGroupID: 10,
			Word:           "アマゾンウェブサービス",
		},
		// Docker関連の表記ゆれ
		{
			ID:             32,
			KeywordGroupID: 11,
			Word:           "Docker",
		},
		{
			ID:             33,
			KeywordGroupID: 11,
			Word:           "docker",
		},
		// Kubernetes関連の表記ゆれ
		{
			ID:             34,
			KeywordGroupID: 12,
			Word:           "Kubernetes",
		},
		{
			ID:             35,
			KeywordGroupID: 12,
			Word:           "K8s",
		},
		{
			ID:             36,
			KeywordGroupID: 12,
			Word:           "k8s",
		},
		// MySQL関連の表記ゆれ
		{
			ID:             37,
			KeywordGroupID: 13,
			Word:           "MySQL",
		},
		{
			ID:             38,
			KeywordGroupID: 13,
			Word:           "mysql",
		},
		// PostgreSQL関連の表記ゆれ
		{
			ID:             39,
			KeywordGroupID: 14,
			Word:           "PostgreSQL",
		},
		{
			ID:             40,
			KeywordGroupID: 14,
			Word:           "postgres",
		},
		{
			ID:             41,
			KeywordGroupID: 14,
			Word:           "Postgres",
		},
		// 設計関連の表記ゆれ
		{
			ID:             42,
			KeywordGroupID: 15,
			Word:           "設計",
		},
		{
			ID:             43,
			KeywordGroupID: 15,
			Word:           "基本設計",
		},
		{
			ID:             44,
			KeywordGroupID: 15,
			Word:           "詳細設計",
		},
		// テスト関連の表記ゆれ
		{
			ID:             45,
			KeywordGroupID: 16,
			Word:           "テスト",
		},
		{
			ID:             46,
			KeywordGroupID: 16,
			Word:           "単体テスト",
		},
		{
			ID:             47,
			KeywordGroupID: 16,
			Word:           "結合テスト",
		},
		// レビュー関連の表記ゆれ
		{
			ID:             48,
			KeywordGroupID: 17,
			Word:           "レビュー",
		},
		{
			ID:             49,
			KeywordGroupID: 17,
			Word:           "コードレビュー",
		},
		{
			ID:             50,
			KeywordGroupID: 17,
			Word:           "設計レビュー",
		},
	}

	for _, keyWord := range keyWords {
		err := tx.Create(&keyWord).Error
		if err != nil {
			return err
		}
	}

	return err
}
