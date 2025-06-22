package application

import (
	cd "business/internal/common/domain"
	"context"
)

// モッククライアント
type mockGmailConnect struct {
	Messages []cd.BasicMessage
	Err      error
}

func (m *mockGmailConnect) GetMessages(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error) {
	return m.Messages, m.Err
}

// func TestGmailConnect_GetMessages(t *testing.T) {
// 	ctx := context.Background()

// 	mockMessages := []cd.BasicMessage{
// 		{
// 			ID:      "abc123",
// 			Subject: "Test Subject",
// 			From:    "from@example.com",
// 			To:      []string{"to@example.com"},
// 			Date:    time.Now(),
// 			Body:    "This is a test body",
// 		},
// 	}

// 	client := &mockGmailConnect{
// 		Messages: mockMessages,
// 		Err:      nil,
// 	}

// 	ei := ei.NewEmailStoreRepository()
// 	ea := ea.new()
// 	conn := New(client)

// 	result, err := conn.GetMessages(ctx, "INBOX", 0)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(result))
// 	assert.Equal(t, "abc123", result[0].ID)
// }
