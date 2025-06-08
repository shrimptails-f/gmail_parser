package application

import (
	"context"
	"testing"
	"time"

	cd "business/internal/common/domain"

	"github.com/stretchr/testify/assert"
)

// モッククライアント
type mockGmailConnect struct {
	Messages []cd.BasicMessage
	Err      error
}

func (m *mockGmailConnect) GetMessages(ctx context.Context, labelName string) ([]cd.BasicMessage, error) {
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

	client := &mockGmailConnect{
		Messages: mockMessages,
		Err:      nil,
	}

	conn := New(client)

	result, err := conn.GetMessages(ctx, "INBOX")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "abc123", result[0].ID)
}
