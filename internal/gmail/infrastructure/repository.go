// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルは認証機能のリポジトリ実装を定義します。
package infrastructure

import (
	"business/internal/gmail/application"
	"business/internal/gmail/domain"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User はデータベースのユーザーテーブルに対応する構造体です
type User struct {
	UserID    uint32    `gorm:"column:user_id;primaryKey;autoIncrement"`
	Email     string    `gorm:"column:email;uniqueIndex;not null"`
	FirstName string    `gorm:"column:first_name;not null"`
	LastName  string    `gorm:"column:last_name;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName はテーブル名を指定します
func (User) TableName() string {
	return "users"
}

// ToDomain はデータベースモデルをドメインモデルに変換します
func (u *User) ToDomain() domain.User {
	return domain.User{
		UserID:    u.UserID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromDomain はドメインモデルをデータベースモデルに変換します
func (u *User) FromDomain(domainUser domain.User) {
	u.UserID = domainUser.UserID
	u.Email = domainUser.Email
	u.FirstName = domainUser.FirstName
	u.LastName = domainUser.LastName
	u.CreatedAt = domainUser.CreatedAt
	u.UpdatedAt = domainUser.UpdatedAt
}

// authRepository は認証機能のリポジトリ実装です
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository は新しい認証リポジトリを作成します
func NewAuthRepository(db *gorm.DB) application.AuthRepository {
	return &authRepository{
		db: db,
	}
}

// GetUserByEmail はメールアドレスでユーザーを取得します
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	// TODO: TDDで実装予定
	return domain.User{}, fmt.Errorf("not implemented")
}

// CreateUser は新しいユーザーを作成します
func (r *authRepository) CreateUser(ctx context.Context, domainUser domain.User) (domain.User, error) {
	// TODO: TDDで実装予定
	return domain.User{}, fmt.Errorf("not implemented")
}
