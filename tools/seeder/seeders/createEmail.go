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
	intPtr := func(i int) *int { return &i }

	emails := []model.Email{
		{
			ID:              "email001",
			Subject:         "【案件】Javaエンジニア募集（リモート可）",
			SenderName:      "田中太郎",
			SenderEmail:     "tanaka@example.com",
			ReceivedDate:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Body:            stringPtr("Java Spring Bootを使用したWebアプリケーション開発案件です。リモートワーク可能です。"),
			Category:        stringPtr("案件"),
			ProjectTitle:    stringPtr("ECサイト構築プロジェクト"),
			EntryTimings:    stringPtr("2024年2月〜"),
			EndTiming:       stringPtr("2024年8月"),
			WorkLocation:    stringPtr("東京都渋谷区（リモート可）"),
			PriceFrom:       intPtr(600000),
			PriceTo:         intPtr(800000),
			RemoteType:      stringPtr("フルリモート可"),
			RemoteFrequency: stringPtr("週5日"),
		},
		{
			ID:              "email002",
			Subject:         "【人材提案】React開発経験者のご紹介",
			SenderName:      "佐藤花子",
			SenderEmail:     "sato@example.com",
			ReceivedDate:    time.Date(2024, 1, 16, 14, 20, 0, 0, time.UTC),
			Body:            stringPtr("React、TypeScriptでの開発経験が豊富なエンジニアをご紹介いたします。"),
			Category:        stringPtr("人材提案"),
			ProjectTitle:    stringPtr("管理画面リニューアル"),
			EntryTimings:    stringPtr("即日〜"),
			EndTiming:       stringPtr("2024年6月"),
			WorkLocation:    stringPtr("東京都新宿区"),
			PriceFrom:       intPtr(700000),
			PriceTo:         intPtr(900000),
			RemoteType:      stringPtr("ハイブリッド"),
			RemoteFrequency: stringPtr("週2-3日出社"),
		},
		{
			ID:              "email003",
			Subject:         "【案件】Python機械学習エンジニア募集",
			SenderName:      "山田次郎",
			SenderEmail:     "yamada@example.com",
			ReceivedDate:    time.Date(2024, 1, 17, 9, 15, 0, 0, time.UTC),
			Body:            stringPtr("Python、機械学習ライブラリを使用したAIシステム開発案件です。"),
			Category:        stringPtr("案件"),
			ProjectTitle:    stringPtr("AIチャットボット開発"),
			EntryTimings:    stringPtr("2024年3月〜"),
			EndTiming:       stringPtr("2024年12月"),
			WorkLocation:    stringPtr("東京都港区"),
			PriceFrom:       intPtr(800000),
			PriceTo:         intPtr(1200000),
			RemoteType:      stringPtr("出社必須"),
			RemoteFrequency: stringPtr("週5日出社"),
		},
		{
			ID:              "email004",
			Subject:         "【案件】Go言語バックエンド開発者募集",
			SenderName:      "鈴木一郎",
			SenderEmail:     "suzuki@example.com",
			ReceivedDate:    time.Date(2024, 1, 18, 16, 45, 0, 0, time.UTC),
			Body:            stringPtr("Go言語でのマイクロサービス開発経験者を募集しています。"),
			Category:        stringPtr("案件"),
			ProjectTitle:    stringPtr("マイクロサービス基盤構築"),
			EntryTimings:    stringPtr("2024年4月〜"),
			EndTiming:       stringPtr("2024年10月"),
			WorkLocation:    stringPtr("東京都千代田区"),
			PriceFrom:       intPtr(750000),
			PriceTo:         intPtr(1000000),
			RemoteType:      stringPtr("フルリモート可"),
			RemoteFrequency: stringPtr("完全リモート"),
		},
		{
			ID:              "email005",
			Subject:         "【人材提案】フルスタックエンジニアのご紹介",
			SenderName:      "高橋美咲",
			SenderEmail:     "takahashi@example.com",
			ReceivedDate:    time.Date(2024, 1, 19, 11, 30, 0, 0, time.UTC),
			Body:            stringPtr("React、Node.js、AWSでの開発経験が豊富なフルスタックエンジニアです。"),
			Category:        stringPtr("人材提案"),
			ProjectTitle:    stringPtr("Webアプリケーション開発"),
			EntryTimings:    stringPtr("2024年2月〜"),
			EndTiming:       stringPtr("長期"),
			WorkLocation:    stringPtr("東京都品川区"),
			PriceFrom:       intPtr(650000),
			PriceTo:         intPtr(850000),
			RemoteType:      stringPtr("ハイブリッド"),
			RemoteFrequency: stringPtr("週1-2日出社"),
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
