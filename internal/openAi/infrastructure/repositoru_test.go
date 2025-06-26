package infrastructure_test

import (
	cd "business/internal/common/domain"
	"business/internal/openAi/infrastructure"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// モック構造体（oa.ClientInterFace のモック）
type mockOpenAIClient struct{}

func (m *mockOpenAIClient) Chat(ctx context.Context, input string) ([]cd.AnalysisResult, error) {
	return []cd.AnalysisResult{}, nil
}

func TestAnalyzer_AnalyzeEmailBody(t *testing.T) {
	mockClient := &mockOpenAIClient{}
	analyzer := infrastructure.New(mockClient)

	ctx := context.Background()
	result, err := analyzer.AnalyzeEmailBody(ctx, "test email content")

	assert.NoError(t, err)
	assert.Equal(t, []cd.AnalysisResult{}, result)
}
