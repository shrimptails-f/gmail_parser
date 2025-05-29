package seeders

import (
	"business/tools/migrations/model"
	"time"

	"gorm.io/gorm"
)

// CreateEmail はメールのサンプルデータを投入する。
func CreateEmail(tx *gorm.DB) error {
	var err error

	// ポインタ型のヘルパー関数
	stringPtr := func(s string) *string { return &s }

	emails := []model.Email{
		{
			ID:           "email001",
			Subject:      "【案件】Javaエンジニア募集（リモート可）",
			SenderName:   "田中太郎",
			SenderEmail:  "tanaka@example.com",
			ReceivedDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Body:         stringPtr("Java Spring Bootを使用したWebアプリケーション開発案件です。リモートワーク可能です。"),
			Category:     "案件",
		},
		{
			ID:           "email002",
			Subject:      "【人材提案】React開発経験者のご紹介",
			SenderName:   "佐藤花子",
			SenderEmail:  "sato@example.com",
			ReceivedDate: time.Date(2024, 1, 16, 14, 20, 0, 0, time.UTC),
			Body:         stringPtr("React、TypeScriptでの開発経験が豊富なエンジニアをご紹介いたします。"),
			Category:     "人材提案",
		},
		{
			ID:           "email003",
			Subject:      "【案件】Python機械学習エンジニア募集",
			SenderName:   "山田次郎",
			SenderEmail:  "yamada@example.com",
			ReceivedDate: time.Date(2024, 1, 17, 9, 15, 0, 0, time.UTC),
			Body:         stringPtr("Python、機械学習ライブラリを使用したAIシステム開発案件です。"),
			Category:     "案件",
		},
		{
			ID:           "email004",
			Subject:      "【案件】Go言語バックエンド開発者募集",
			SenderName:   "鈴木一郎",
			SenderEmail:  "suzuki@example.com",
			ReceivedDate: time.Date(2024, 1, 18, 16, 45, 0, 0, time.UTC),
			Body:         stringPtr("Go言語でのマイクロサービス開発経験者を募集しています。"),
			Category:     "案件",
		},
		{
			ID:           "email005",
			Subject:      "【人材提案】フルスタックエンジニアのご紹介",
			SenderName:   "高橋美咲",
			SenderEmail:  "takahashi@example.com",
			ReceivedDate: time.Date(2024, 1, 19, 11, 30, 0, 0, time.UTC),
			Body:         stringPtr("React、Node.js、AWSでの開発経験が豊富なフルスタックエンジニアです。"),
			Category:     "人材提案",
		},
	}

	for _, email := range emails {
		err := tx.Create(&email).Error
		if err != nil {
			return err
		}
	}

	return err
}
