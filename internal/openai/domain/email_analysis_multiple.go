// Package domain はAI機能のドメイン層を提供します。
// このファイルは複数案件対応のメール分析に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"time"
)

// EmailAnalysisMultipleResult は複数案件対応のメール分析結果のドメインモデルです
type EmailAnalysisMultipleResult struct {
	// Gmailの情報
	GmailID   string    `json:"gmail_id"`
	Subject   string    `json:"subject"`
	From      string    `json:"from"`
	FromEmail string    `json:"from_email"`
	Date      time.Time `json:"date"`
	Body      string    `json:"body"`

	// AI分析結果（複数案件対応）
	MailCategory string                    `json:"メール区分"`
	Projects     []ProjectAnalysisResult   `json:"案件一覧"`
	Candidates   []CandidateAnalysisResult `json:"人材一覧"`
}

// ProjectAnalysisResult は案件分析結果のドメインモデルです
type ProjectAnalysisResult struct {
	ProjectName         string   `json:"案件名"`
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

// CandidateAnalysisResult は人材分析結果のドメインモデルです
type CandidateAnalysisResult struct {
	CandidateName    string `json:"人材名"`
	ExperienceYears  *int   `json:"経験年数"`
	SkillsSummary    string `json:"スキルまとめ"`
	AvailabilityDate string `json:"参画可能日"`
}

// 複数案件対応のメール分析関連のドメインエラー
var (
	ErrNoProjectsFound      = errors.New("案件情報が見つかりません")
	ErrNoCandidatesFound    = errors.New("人材情報が見つかりません")
	ErrInvalidProjectData   = errors.New("無効な案件データです")
	ErrInvalidCandidateData = errors.New("無効な人材データです")
)

// NewEmailAnalysisMultipleResult は複数案件対応のメール分析結果を作成します
func NewEmailAnalysisMultipleResult(id, subject, from, fromEmail, body string, date time.Time) *EmailAnalysisMultipleResult {
	return &EmailAnalysisMultipleResult{
		GmailID:    id,
		Subject:    subject,
		From:       from,
		FromEmail:  fromEmail,
		Date:       date,
		Body:       body,
		Projects:   make([]ProjectAnalysisResult, 0),
		Candidates: make([]CandidateAnalysisResult, 0),
	}
}

// AddProject は案件分析結果を追加します
func (r *EmailAnalysisMultipleResult) AddProject(project ProjectAnalysisResult) {
	r.Projects = append(r.Projects, project)
}

// AddCandidate は人材分析結果を追加します
func (r *EmailAnalysisMultipleResult) AddCandidate(candidate CandidateAnalysisResult) {
	r.Candidates = append(r.Candidates, candidate)
}

// HasProjects は案件情報が存在するかチェックします
func (r *EmailAnalysisMultipleResult) HasProjects() bool {
	return len(r.Projects) > 0
}

// HasCandidates は人材情報が存在するかチェックします
func (r *EmailAnalysisMultipleResult) HasCandidates() bool {
	return len(r.Candidates) > 0
}

// IsValid は複数案件対応のメール分析結果の妥当性をチェックします
func (r *EmailAnalysisMultipleResult) IsValid() error {
	if r.GmailID == "" {
		return errors.New("メールIDが空です")
	}
	if r.Body == "" {
		return errors.New("本文が空です")
	}
	if r.MailCategory == "" {
		return errors.New("メール区分が空です")
	}

	// 案件メールの場合は案件情報または人材情報が必要
	if r.MailCategory == "案件" && !r.HasProjects() && !r.HasCandidates() {
		return errors.New("案件メールには案件情報または人材情報が必要です")
	}

	return nil
}

// GetProjectCount は案件数を返します
func (r *EmailAnalysisMultipleResult) GetProjectCount() int {
	return len(r.Projects)
}

// GetCandidateCount は人材数を返します
func (r *EmailAnalysisMultipleResult) GetCandidateCount() int {
	return len(r.Candidates)
}
