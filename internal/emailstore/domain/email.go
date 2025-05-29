// Package domain はメール保存機能のドメイン層を提供します。
// このファイルはメール保存に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"time"
)

// Email は全メール共通の基本情報を表すドメインモデルです
type Email struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Subject   string    `gorm:"type:text;not null" json:"subject"`
	From      string    `gorm:"type:varchar(500);not null" json:"from"`
	FromEmail string    `gorm:"type:varchar(255);not null" json:"from_email"`
	Date      time.Time `gorm:"not null" json:"date"`
	Body      string    `gorm:"type:longtext;not null" json:"body"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// EmailProject は案件メール専用の詳細情報を表すドメインモデルです
type EmailProject struct {
	ID                  uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailID             string  `gorm:"type:varchar(255);not null;index" json:"email_id"`
	MailCategory        string  `gorm:"type:varchar(100)" json:"mail_category"`
	EndPeriod           string  `gorm:"type:varchar(100)" json:"end_period"`
	WorkLocation        string  `gorm:"type:varchar(500)" json:"work_location"`
	PriceFrom           *int    `gorm:"type:int" json:"price_from"`
	PriceTo             *int    `gorm:"type:int" json:"price_to"`
	RemoteWorkCategory  string  `gorm:"type:varchar(100)" json:"remote_work_category"`
	RemoteWorkFrequency *string `gorm:"type:varchar(100)" json:"remote_work_frequency"`
	// 一覧画面用のカンマ区切り文字列（二重管理）
	TechnologiesText string    `gorm:"type:text" json:"technologies_text"`
	PositionsText    string    `gorm:"type:text" json:"positions_text"`
	WorkTypesText    string    `gorm:"type:text" json:"work_types_text"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	Email        Email         `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	EntryTimings []EntryTiming `gorm:"foreignKey:EmailProjectID" json:"entry_timings"`
}

// EntryTiming は案件の入場時期を正規化管理するドメインモデルです
type EntryTiming struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailProjectID uint      `gorm:"not null;index" json:"email_project_id"`
	Timing         string    `gorm:"type:varchar(100);not null" json:"timing"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	EmailProject EmailProject `gorm:"foreignKey:EmailProjectID" json:"email_project"`
}

// KeywordGroup は正規化された技術キーワードのマスタを表すドメインモデルです
type KeywordGroup struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	KeyWords           []KeyWord           `gorm:"foreignKey:KeywordGroupID" json:"key_words"`
	EmailKeywordGroups []EmailKeywordGroup `gorm:"foreignKey:KeywordGroupID" json:"email_keyword_groups"`
}

// KeyWord はキーワードの表記ゆれをKeywordGroupに紐付けるドメインモデルです
type KeyWord struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	KeywordGroupID uint      `gorm:"not null;index" json:"keyword_group_id"`
	Word           string    `gorm:"type:varchar(100);not null" json:"word"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	KeywordGroup KeywordGroup `gorm:"foreignKey:KeywordGroupID" json:"keyword_group"`
}

// EmailKeywordGroup はEmailとKeywordGroupの多対多中間テーブルを表すドメインモデルです
type EmailKeywordGroup struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailID        string    `gorm:"type:varchar(255);not null;index" json:"email_id"`
	KeywordGroupID uint      `gorm:"not null;index" json:"keyword_group_id"`
	Type           string    `gorm:"type:varchar(50);not null" json:"type"` // language, framework, skill_must, skill_want
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	Email        Email        `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	KeywordGroup KeywordGroup `gorm:"foreignKey:KeywordGroupID" json:"keyword_group"`
}

// PositionGroup は正規化されたポジション名のマスタを表すドメインモデルです
type PositionGroup struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	PositionWords       []PositionWord       `gorm:"foreignKey:PositionGroupID" json:"position_words"`
	EmailPositionGroups []EmailPositionGroup `gorm:"foreignKey:PositionGroupID" json:"email_position_groups"`
}

// PositionWord はポジションの表記ゆれをPositionGroupに紐付けるドメインモデルです
type PositionWord struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PositionGroupID uint      `gorm:"not null;index" json:"position_group_id"`
	Word            string    `gorm:"type:varchar(100);not null" json:"word"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	PositionGroup PositionGroup `gorm:"foreignKey:PositionGroupID" json:"position_group"`
}

// EmailPositionGroup はEmailとPositionGroupの多対多中間テーブルを表すドメインモデルです
type EmailPositionGroup struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailID         string    `gorm:"type:varchar(255);not null;index" json:"email_id"`
	PositionGroupID uint      `gorm:"not null;index" json:"position_group_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	Email         Email         `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	PositionGroup PositionGroup `gorm:"foreignKey:PositionGroupID" json:"position_group"`
}

// WorkTypeGroup は正規化された業務種別マスタを表すドメインモデルです
type WorkTypeGroup struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	WorkTypeWords       []WorkTypeWord       `gorm:"foreignKey:WorkTypeGroupID" json:"work_type_words"`
	EmailWorkTypeGroups []EmailWorkTypeGroup `gorm:"foreignKey:WorkTypeGroupID" json:"email_work_type_groups"`
}

// WorkTypeWord は業務表記ゆれをWorkTypeGroupに紐付けるドメインモデルです
type WorkTypeWord struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkTypeGroupID uint      `gorm:"not null;index" json:"work_type_group_id"`
	Word            string    `gorm:"type:varchar(100);not null" json:"word"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	WorkTypeGroup WorkTypeGroup `gorm:"foreignKey:WorkTypeGroupID" json:"work_type_group"`
}

// EmailWorkTypeGroup はEmailとWorkTypeGroupの多対多中間テーブルを表すドメインモデルです
type EmailWorkTypeGroup struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmailID         string    `gorm:"type:varchar(255);not null;index" json:"email_id"`
	WorkTypeGroupID uint      `gorm:"not null;index" json:"work_type_group_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// リレーション
	Email         Email         `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	WorkTypeGroup WorkTypeGroup `gorm:"foreignKey:WorkTypeGroupID" json:"work_type_group"`
}

// ドメインエラー
var (
	ErrEmailNotFound      = errors.New("メールが見つかりません")
	ErrEmailAlreadyExists = errors.New("メールが既に存在します")
	ErrInvalidEmailData   = errors.New("無効なメールデータです")
)

// TableName はテーブル名を指定します
func (Email) TableName() string {
	return "emails"
}

func (EmailProject) TableName() string {
	return "email_projects"
}

func (EntryTiming) TableName() string {
	return "entry_timings"
}

func (KeywordGroup) TableName() string {
	return "keyword_groups"
}

func (KeyWord) TableName() string {
	return "key_words"
}

func (EmailKeywordGroup) TableName() string {
	return "email_keyword_groups"
}

func (PositionGroup) TableName() string {
	return "position_groups"
}

func (PositionWord) TableName() string {
	return "position_words"
}

func (EmailPositionGroup) TableName() string {
	return "email_position_groups"
}

func (WorkTypeGroup) TableName() string {
	return "work_type_groups"
}

func (WorkTypeWord) TableName() string {
	return "work_type_words"
}

func (EmailWorkTypeGroup) TableName() string {
	return "email_work_type_groups"
}
