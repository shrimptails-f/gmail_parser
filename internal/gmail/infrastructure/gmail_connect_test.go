package infrastructure

import (
	"context"
	"testing"
	"time"

	cd "business/internal/common/domain"

	"github.com/stretchr/testify/assert"
)

// モッククライアント
type mockClient struct {
	Messages []cd.BasicMessage
	Err      error
}

func (m *mockClient) ListMessageIDs(ctx context.Context, max int64) ([]string, error) {
	return nil, nil
}

func (m *mockClient) GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error) {
	return m.Messages, m.Err
}

func TestGmailConnect_GetMessages(t *testing.T) {
	ctx := context.Background()

	mockMessages := []cd.BasicMessage{
		{
			ID:      "abc123",
			Subject: "Test Subject",
			From:    "from@example.com",
			To:      []string{"to@example.com"},
			Date:    time.Now(),
			Body:    "This is a test body",
		},
	}

	client := &mockClient{
		Messages: mockMessages,
		Err:      nil,
	}
	conn := New(client)

	result, err := conn.GetMessages(ctx, "INBOX", 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "abc123", result[0].ID)
}
