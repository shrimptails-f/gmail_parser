package seeders

import (
	"business/tools/migrations/model"

	"gorm.io/gorm"
)

// CreatePosition はポジションのサンプルデータを投入する。
func CreatePosition(tx *gorm.DB) error {
	var err error
	positions := []model.Position{
		{
			ID:   1,
			Name: "PG",
		},
		{
			ID:   2,
			Name: "SE",
		},
		{
			ID:   3,
			Name: "PL",
		},
		{
			ID:   4,
			Name: "PM",
		},
		{
			ID:   5,
			Name: "PMO",
		},
		{
			ID:   6,
			Name: "アーキテクト",
		},
		{
			ID:   7,
			Name: "テックリード",
		},
		{
			ID:   8,
			Name: "フロントエンドエンジニア",
		},
		{
			ID:   9,
			Name: "バックエンドエンジニア",
		},
		{
			ID:   10,
			Name: "フルスタックエンジニア",
		},
		{
			ID:   11,
			Name: "インフラエンジニア",
		},
		{
			ID:   12,
			Name: "DevOpsエンジニア",
		},
		{
			ID:   13,
			Name: "QAエンジニア",
		},
		{
			ID:   14,
			Name: "データエンジニア",
		},
		{
			ID:   15,
			Name: "機械学習エンジニア",
		},
	}

	for _, position := range positions {
		err := tx.Create(&position).Error
		if err != nil {
			return err
		}
	}

	return err
}
