package main

// EmailAnalysisResult はメール分析結果のドメインモデルです
type EmailAnalysisResult struct {
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
