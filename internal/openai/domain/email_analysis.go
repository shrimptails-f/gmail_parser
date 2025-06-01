// Package domain はAI機能のドメイン層を提供します。
// このファイルはメール分析に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"time"
)

// EmailAnalysisResult はメール分析結果のドメインモデルです
type EmailAnalysisResult struct {
	// Gmailの情報
	GmailID   string    `json:"gmail_id"`
	Subject   string    `json:"subject"`
	From      string    `json:"from"`
	FromEmail string    `json:"from_email"`
	Date      time.Time `json:"date"`
	Body      string    `json:"body"`

	// AI分析結果
	MailCategory        string   `json:"メール区分"`
	StartPeriod         []string `json:"入場時期・開始時期"`
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
	RemoteWorkCategory  string   `json:"リモートワーク区分"`
	RemoteWorkFrequency *string  `json:"リモートワークの頻度"`
}

// EmailAnalysisRequest はメール分析リクエストのドメインモデルです
type EmailAnalysisRequest struct {
	EmailText string            `json:"email_text"`
	MessageID string            `json:"message_id"`
	Subject   string            `json:"subject"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// メール分析関連のドメインエラー
var (
	ErrEmptyEmailText        = errors.New("分析対象のメール本文が空です")
	ErrEmailAnalysisTimeout  = errors.New("メール分析がタイムアウトしました")
	ErrEmailAnalysisAPIError = errors.New("メール分析APIでエラーが発生しました")
	ErrInvalidEmailFormat    = errors.New("無効なメール形式です")
)

// NewEmailAnalysisRequest はメール分析リクエストを作成します
func NewEmailAnalysisRequest(emailText, messageID, subject string) *EmailAnalysisRequest {
	return &EmailAnalysisRequest{
		EmailText: emailText,
		MessageID: messageID,
		Subject:   subject,
		Metadata:  make(map[string]string),
	}
}

// IsValid はメール分析リクエストの妥当性をチェックします
func (r *EmailAnalysisRequest) IsValid() error {
	if r.EmailText == "" {
		return ErrEmptyEmailText
	}
	if len(r.EmailText) > 100000 { // 100KB制限
		return errors.New("メール本文が長すぎます（最大100KB）")
	}
	return nil
}
