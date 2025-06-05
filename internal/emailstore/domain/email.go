// Package domain はメール保存機能のドメイン層を提供します。
// このファイルはメール保存に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"strings"
	"time"
)

// AnalysisResult は全メール共通の基本情報を表すドメインモデルです
type AnalysisResult struct {
	GmailID      string    `json:"gmail_id"`
	ReceivedDate time.Time `json:"received_date"`
	Summary      string    `json:"summary"`
	Subject      string    `json:"subject"`
	From         string    `json:"from"`
	FromEmail    string    `json:"from_email"`
	Body         string    `json:"body"`

	Category            string   `json:"メール区分"`
	ProjectName         string   `json:"案件名"`
	StartPeriod         []string `json:"開始時期"`
	EndPeriod           string   `json:"終了時期"`
	WorkLocation        string   `json:"勤務場所"`
	PriceFrom           *int     `json:"単価FROM"`
	PriceTo             *int     `json:"単価TO"`
	Languages           []string `json:"言語"`
	Frameworks          []string `json:"フレームワーク"`
	Positions           []string `json:"ポジション"`
	WorkTypes           []string `json:"業務"`
	RequiredSkillsMust  []string `json:"求めるスキル MUST"`
	RequiredSkillsWant  []string `json:"求めるスキル WANT"`
	RemoteWorkCategory  *string  `json:"リモートワーク区分"`
	RemoteWorkFrequency *string  `json:"リモートワークの頻度"`
}

// SenderName は From フィールドから送信者名を抽出します
func (a *AnalysisResult) SenderName() string {
	if idx := strings.Index(a.From, "<"); idx > 0 {
		return strings.TrimSpace(a.From[:idx])
	}
	return a.From
}

// SenderEmail は From フィールドからメールアドレスを抽出します
func (a *AnalysisResult) SenderEmail() string {
	start := strings.Index(a.From, "<")
	end := strings.Index(a.From, ">")
	if start >= 0 && end > start {
		return a.From[start+1 : end]
	}
	return a.From
}

// Email は全メール共通の基本情報を表すドメインモデルです
type Email struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`           // オートインクリメントID
	GmailID      string    `gorm:"size:32;index"`                      // GメールID
	Subject      string    `gorm:"type:text;not null" json:"subject"`  // 件名
	SenderName   string    `gorm:"size:255" json:"sender_name"`        // 差出人名
	SenderEmail  string    `gorm:"size:255;index" json:"sender_email"` // メールアドレス
	ReceivedDate time.Time `gorm:"index" json:"received_date"`         // 受信日
	Body         *string   `gorm:"type:longtext" json:"body"`          // 本文
	Category     string    `gorm:"size:50;index" json:"category"`      // 種別（案件 / 人材提案）
	CreatedAt    time.Time `json:"created_at"`                         // 作成日時
	UpdatedAt    time.Time `json:"updated_at"`                         // 更新日時

	// 子テーブル
	EmailProject        *EmailProject        `gorm:"foreignKey:EmailID;references:ID" json:"email_project"`          // 案件情報（1対1）
	EmailCandidate      *EmailCandidate      `gorm:"foreignKey:EmailID;references:ID" json:"email_candidate"`        // 人材情報（1対1）
	EmailKeywordGroups  []EmailKeywordGroup  `gorm:"foreignKey:EmailID;references:ID" json:"email_keyword_groups"`   // 技術キーワード（1対多）
	EmailPositionGroups []EmailPositionGroup `gorm:"foreignKey:EmailID;references:ID" json:"email_position_groups"`  // ポジション（1対多）
	EmailWorkTypeGroups []EmailWorkTypeGroup `gorm:"foreignKey:EmailID;references:ID" json:"email_work_type_groups"` // 業務内容（1対多）
}

// EmailProject は案件メール専用の詳細情報を表すドメインモデルです
type EmailProject struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`      // オートインクリメントID
	EmailID      uint    `gorm:"index"`                         // メールID（emails.idと同じ）
	ProjectTitle *string `gorm:"size:255" json:"project_title"` // 案件名

	// 表示用（カンマ区切り）
	EntryTiming *string `gorm:"type:text" json:"entry_timing"` // 入場時期（"2025/06/01,2025/07/01"）
	Languages   *string `gorm:"type:text" json:"languages"`    // 言語（"PHP,TypeScript"）
	Frameworks  *string `gorm:"type:text" json:"frameworks"`   // フレームワーク（"React,Laravel"）
	Positions   *string `gorm:"type:text" json:"positions"`    // ポジション（"PM,SE"）
	WorkTypes   *string `gorm:"type:text" json:"work_types"`   // 業務内容（"バックエンド実装,インフラ構築"）
	MustSkills  *string `gorm:"type:text" json:"must_skills"`  // MUSTスキル（"CMS知識,PowerCMS"）
	WantSkills  *string `gorm:"type:text" json:"want_skills"`  // WANTスキル（"MT,Adobe製品経験"）

	// その他項目
	EndTiming       *string   `gorm:"size:255" json:"end_timing"`          // 終了時期
	WorkLocation    *string   `gorm:"size:255;index" json:"work_location"` // 勤務場所
	PriceFrom       *int      `gorm:"type:int" json:"price_from"`          // 単価FROM
	PriceTo         *int      `gorm:"type:int" json:"price_to"`            // 単価TO
	RemoteType      *string   `gorm:"size:50" json:"remote_type"`          // リモート区分
	RemoteFrequency *string   `gorm:"size:255" json:"remote_frequency"`    // リモート頻度
	CreatedAt       time.Time `json:"created_at"`                          // 作成日時
	UpdatedAt       time.Time `json:"updated_at"`                          // 更新日時

	// リレーション
	Email        Email         `gorm:"foreignKey:EmailID;references:ID" json:"email"`  // 親メール
	EntryTimings []EntryTiming `gorm:"foreignKey:EmailProjectID" json:"entry_timings"` // 入場時期（1対多）
}

// EmailCandidate は人材メール専用の詳細情報を表すドメインモデルです（将来拡張用）
type EmailCandidate struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // オートインクリメントID
	EmailID   uint      `gorm:"index"`                    // メールID（emails.idと同じ）
	CreatedAt time.Time `json:"created_at"`               // 作成日時
	UpdatedAt time.Time `json:"updated_at"`               // 更新日時

	// リレーション
	Email Email `gorm:"foreignKey:EmailID;references:ID" json:"email"`
}

// EntryTiming は案件の入場時期を正規化管理するドメインモデルです
type EntryTiming struct {
	ID             uint      `gorm:"primaryKey" json:"id"`                  // ID
	EmailProjectID uint      `gorm:"size:32;index" json:"email_project_id"` // 紐づく案件メールID
	StartDate      string    `gorm:"size:20;not null" json:"start_date"`    // 入場日（例: "2025/06/01"）
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// KeywordGroup は正規化された技術キーワードのマスタを表すドメインモデルです
type KeywordGroup struct {
	KeywordGroupID uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"unique;size:255;not null"`
	Type           string `gorm:"type:enum('language','framework','must','want','other');not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Words []KeyWord `gorm:"foreignKey:KeywordGroupID;references:KeywordGroupID" json:"words"` // 表記ゆれ一覧（統合テスト時はコメントアウト）
}

// 登録単語1に対してKeywordGroupを複数登録するための中間テーブル
type KeywordGroupWordLink struct {
	KeywordGroupID uint `gorm:"primaryKey"`
	KeyWordID      uint `gorm:"primaryKey"`
	CreatedAt      time.Time

	// 循環してて完全に積んでるのでコメントアウト
	// KeyWord KeyWord `gorm:"foreignKey:KeyWordID;references:ID"`
}

// KeyWord はキーワードの表記ゆれをKeywordGroupに紐付けるドメインモデルです
type KeyWord struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Word      string `gorm:"size:255;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EmailKeywordGroup はEmailとKeywordGroupの多対多中間テーブルを表すドメインモデルです
type EmailKeywordGroup struct {
	EmailID          uint
	KeywordGroupID   uint
	KeywordGroupName string // キーワードグループの単語名
	KeywordWord      string // キーワードの単語名

	Email Email `gorm:"foreignKey:EmailID;references:ID" json:"email"`
}

// PositionGroup は正規化されたポジション名のマスタを表すドメインモデルです
type PositionGroup struct {
	PositionGroupID uint      `gorm:"primaryKey" json:"position_group_id"`  // ポジショングループID
	Name            string    `gorm:"unique;size:100;not null" json:"name"` // 正規化されたポジション名（例: "PM"）
	CreatedAt       time.Time `json:"created_at"`                           // 作成日時
	UpdatedAt       time.Time `json:"updated_at"`                           // 更新日時

	// Words []PositionWord `gorm:"foreignKey:PositionGroupID;references:PositionGroupID" json:"words"` // 表記ゆれ一覧（統合テスト時はコメントアウト）
}

// PositionWord はポジションの表記ゆれをPositionGroupに紐付けるドメインモデルです
type PositionWord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`              // 表記ID
	PositionGroupID uint      `gorm:"not null" json:"position_group_id"` // 紐づくポジショングループID
	Word            string    `gorm:"size:100;not null" json:"word"`     // 表記（例: "Project Manager", "ＰＭ"）
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EmailPositionGroup はEmailとPositionGroupの多対多中間テーブルを表すドメインモデルです
type EmailPositionGroup struct {
	EmailID         uint `gorm:"primaryKey;size:32" json:"email_id"`  // メールID
	PositionGroupID uint `gorm:"primaryKey" json:"position_group_id"` // ポジショングループID

	// リレーション
	Email Email `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	// PositionGroup PositionGroup `gorm:"foreignKey:PositionGroupID;references:PositionGroupID" json:"position_group"` // 統合テスト時はコメントアウト
}

// WorkTypeGroup は正規化された業務種別マスタを表すドメインモデルです
type WorkTypeGroup struct {
	WorkTypeGroupID uint      `gorm:"primaryKey" json:"work_type_group_id"` // 業務グループID
	Name            string    `gorm:"unique;size:100;not null" json:"name"` // 正規化された業務名（例: "バックエンド開発"）
	CreatedAt       time.Time `json:"created_at"`                           // 作成日時
	UpdatedAt       time.Time `json:"updated_at"`                           // 更新日時

	// Words []WorkTypeWord `gorm:"foreignKey:WorkTypeGroupID;references:WorkTypeGroupID" json:"words"` // 表記ゆれ一覧（統合テスト時はコメントアウト）
}

// WorkTypeWord は業務表記ゆれをWorkTypeGroupに紐付けるドメインモデルです
type WorkTypeWord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`               // 表記ID
	WorkTypeGroupID uint      `gorm:"not null" json:"work_type_group_id"` // 紐づく業務グループID
	Word            string    `gorm:"size:100;not null" json:"word"`      // 表記（例: "BE実装", "バックエンド構築"）
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EmailWorkTypeGroup はEmailとWorkTypeGroupの多対多中間テーブルを表すドメインモデルです
type EmailWorkTypeGroup struct {
	EmailID         uint `gorm:"primaryKey;size:32" json:"email_id"`   // メールID
	WorkTypeGroupID uint `gorm:"primaryKey" json:"work_type_group_id"` // 業務グループID

	// リレーション
	Email Email `gorm:"foreignKey:EmailID;references:ID" json:"email"`
	// WorkTypeGroup WorkTypeGroup `gorm:"foreignKey:WorkTypeGroupID;references:WorkTypeGroupID" json:"work_type_group"` // 統合テスト時はコメントアウト
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

func (EmailCandidate) TableName() string {
	return "email_candidates"
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
