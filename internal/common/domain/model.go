package domain

import (
	"strings"
	"time"
)

// AnalysisResult は全メール共通の基本情報を表すドメインモデルです
type AnalysisResult struct {
	MailCategory        string   `json:"メール区分"`
	ProjectTitle        string   `json:"案件名"`
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

// Email は全メール共通の基本情報を表すドメインモデルです
type Email struct {
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
func (e *Email) SenderName() string {
	if idx := strings.Index(e.From, "<"); idx > 0 {
		return strings.TrimSpace(e.From[:idx])
	}
	return e.From
}

// SenderEmail は From フィールドからメールアドレスを抽出します
func (e *Email) SenderEmail() string {
	start := strings.Index(e.From, "<")
	end := strings.Index(e.From, ">")
	if start >= 0 && end > start {
		return e.From[start+1 : end]
	}
	return e.From
}
