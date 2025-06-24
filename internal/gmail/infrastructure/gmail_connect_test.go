package infrastructure

import (
	cd "business/internal/common/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// モッククライアント
type mockClient struct {
	Ids     []string
	Message cd.BasicMessage
	Err     error
}

func (m *mockClient) ListMessageIDs(ctx context.Context, max int64) ([]string, error) {
	return nil, nil
}

func (m *mockClient) GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	return m.Ids, m.Err
}

func (m *mockClient) GetGmailDetail(id string) (cd.BasicMessage, error) {
	return m.Message, m.Err
}

func TestGmailConnect_GetMessages(t *testing.T) {
	ctx := context.Background()
	client := &mockClient{
		Ids: []string{"abc123"},
		Err: nil,
	}
	conn := New(client)

	result, err := conn.gc.GetMessagesByLabelName(ctx, "INBOX", 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "abc123", result[0])

}

func TestGmailConnect_GetGmailDetail(t *testing.T) {
	mockMessages := cd.BasicMessage{
		ID:      "abc123",
		Subject: "Test Subject",
		From:    "from@example.com",
		To:      []string{"to@example.com"},
		Date:    time.Now(),
		Body:    "This is a test body",
	}

	client := &mockClient{
		Message: mockMessages,
		Err:     nil,
	}
	conn := New(client)

	result, err := conn.gc.GetGmailDetail("abc123")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", result.ID)
}
