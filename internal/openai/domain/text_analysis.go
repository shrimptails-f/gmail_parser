// Package domain はAI機能のドメイン層を提供します。
// このファイルはテキスト字句解析に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"time"
)

// TextAnalysisResult はテキスト字句解析結果のドメインモデルです
type TextAnalysisResult struct {
	MessageID   string                 `json:"message_id"`
	Subject     string                 `json:"subject"`
	From        string                 `json:"from"`
	AnalyzedAt  time.Time              `json:"analyzed_at"`
	Sentiment   SentimentAnalysis      `json:"sentiment"`
	Keywords    []Keyword              `json:"keywords"`
	Entities    []Entity               `json:"entities"`
	Summary     string                 `json:"summary"`
	Categories  []Category             `json:"categories"`
	Language    string                 `json:"language"`
	Confidence  float64                `json:"confidence"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

// SentimentAnalysis は感情分析結果のドメインモデルです
type SentimentAnalysis struct {
	Score      float64 `json:"score"`      // -1.0 (negative) to 1.0 (positive)
	Magnitude  float64 `json:"magnitude"`  // 0.0 to infinity
	Label      string  `json:"label"`      // POSITIVE, NEGATIVE, NEUTRAL, MIXED
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
}

// Keyword はキーワード抽出結果のドメインモデルです
type Keyword struct {
	Text      string  `json:"text"`
	Relevance float64 `json:"relevance"` // 0.0 to 1.0
	Count     int     `json:"count"`
	Category  string  `json:"category"`
}

// Entity はエンティティ抽出結果のドメインモデルです
type Entity struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`       // PERSON, ORGANIZATION, LOCATION, etc.
	Salience   float64         `json:"salience"`   // 0.0 to 1.0
	Confidence float64         `json:"confidence"` // 0.0 to 1.0
	Mentions   []EntityMention `json:"mentions"`
}

// EntityMention はエンティティの言及箇所のドメインモデルです
type EntityMention struct {
	Text   string `json:"text"`
	Type   string `json:"type"`   // PROPER, COMMON
	Offset int    `json:"offset"` // テキスト内の位置
	Length int    `json:"length"` // 言及の長さ
}

// Category はカテゴリ分類結果のドメインモデルです
type Category struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
	Path       string  `json:"path"`       // カテゴリの階層パス
}

// TextAnalysisRequest はテキスト字句解析リクエストのドメインモデルです
type TextAnalysisRequest struct {
	Text     string            `json:"text"`
	Language string            `json:"language,omitempty"` // 言語コード（例：ja, en）
	Options  AnalysisOptions   `json:"options"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// AnalysisOptions は字句解析のオプション設定です
type AnalysisOptions struct {
	EnableSentiment  bool `json:"enable_sentiment"`
	EnableKeywords   bool `json:"enable_keywords"`
	EnableEntities   bool `json:"enable_entities"`
	EnableSummary    bool `json:"enable_summary"`
	EnableCategories bool `json:"enable_categories"`
	MaxKeywords      int  `json:"max_keywords"`
	MaxSummaryLength int  `json:"max_summary_length"`
}

// テキスト字句解析関連のドメインエラー
var (
	ErrEmptyText             = errors.New("解析対象のテキストが空です")
	ErrInvalidLanguage       = errors.New("無効な言語コードです")
	ErrAnalysisTimeout       = errors.New("字句解析がタイムアウトしました")
	ErrAnalysisAPIError      = errors.New("字句解析APIでエラーが発生しました")
	ErrInvalidAnalysisType   = errors.New("無効な解析タイプです")
	ErrAnalysisQuotaExceeded = errors.New("字句解析APIのクォータを超過しました")
)

// NewTextAnalysisRequest はテキスト字句解析リクエストを作成します
func NewTextAnalysisRequest(text string) *TextAnalysisRequest {
	return &TextAnalysisRequest{
		Text: text,
		Options: AnalysisOptions{
			EnableSentiment:  true,
			EnableKeywords:   true,
			EnableEntities:   true,
			EnableSummary:    true,
			EnableCategories: true,
			MaxKeywords:      10,
			MaxSummaryLength: 200,
		},
		Metadata: make(map[string]string),
	}
}

// IsValid はテキスト字句解析リクエストの妥当性をチェックします
func (r *TextAnalysisRequest) IsValid() error {
	if r.Text == "" {
		return ErrEmptyText
	}
	if len(r.Text) > 100000 { // 100KB制限
		return errors.New("テキストが長すぎます（最大100KB）")
	}
	return nil
}

// IsPositive は感情分析結果がポジティブかどうかを判定します
func (s *SentimentAnalysis) IsPositive() bool {
	return s.Score > 0.1 && s.Label == "POSITIVE"
}

// IsNegative は感情分析結果がネガティブかどうかを判定します
func (s *SentimentAnalysis) IsNegative() bool {
	return s.Score < -0.1 && s.Label == "NEGATIVE"
}

// IsNeutral は感情分析結果がニュートラルかどうかを判定します
func (s *SentimentAnalysis) IsNeutral() bool {
	return s.Score >= -0.1 && s.Score <= 0.1 && s.Label == "NEUTRAL"
}

// GetHighConfidenceKeywords は高い信頼度のキーワードのみを返します
func (r *TextAnalysisResult) GetHighConfidenceKeywords(threshold float64) []Keyword {
	var result []Keyword
	for _, keyword := range r.Keywords {
		if keyword.Relevance >= threshold {
			result = append(result, keyword)
		}
	}
	return result
}

// GetEntitiesByType は指定されたタイプのエンティティのみを返します
func (r *TextAnalysisResult) GetEntitiesByType(entityType string) []Entity {
	var result []Entity
	for _, entity := range r.Entities {
		if entity.Type == entityType {
			result = append(result, entity)
		}
	}
	return result
}
