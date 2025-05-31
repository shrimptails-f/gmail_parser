package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/gmail/v1"
)

func TestGmailMessageService_convertToGmailMessage(t *testing.T) {
	service := &gmailMessageService{}

	tests := []struct {
		name     string
		message  *gmail.Message
		expected string // 期待される本文（HTMLタグが除去された状態）
	}{
		{
			name: "プレーンテキストの場合にそのまま返すこと",
			message: &gmail.Message{
				Id: "test-id-1",
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "Subject", Value: "テストメール"},
						{Name: "From", Value: "test@example.com"},
						{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0900"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/plain",
							Body: &gmail.MessagePartBody{
								Data: "SGVsbG8sIFdvcmxkIQ==", // "Hello, World!" in base64
							},
						},
					},
				},
			},
			expected: "Hello, World!",
		},
		{
			name: "HTMLメールの場合にHTMLタグを除去すること",
			message: &gmail.Message{
				Id: "test-id-2",
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "Subject", Value: "HTMLテストメール"},
						{Name: "From", Value: "test@example.com"},
						{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0900"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/html",
							Body: &gmail.MessagePartBody{
								Data: "PGRpdj5IZWxsbywgPHN0cm9uZz5Xb3JsZDwvc3Ryb25nPiE8L2Rpdj4=", // "<div>Hello, <strong>World</strong>!</div>" in base64
							},
						},
					},
				},
			},
			expected: "Hello, World!",
		},
		{
			name: "スタイル属性付きHTMLの場合にHTMLタグを除去すること",
			message: &gmail.Message{
				Id: "test-id-3",
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "Subject", Value: "スタイル付きHTMLテストメール"},
						{Name: "From", Value: "test@example.com"},
						{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0900"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/html",
							Body: &gmail.MessagePartBody{
								Data: "PGRpdiBzdHlsZT0iY29sb3I6IHJlZDsiPjxmb250IHN0eWxlPSJmb250LXNpemU6IDEycHg7Ij5IZWxsbywgV29ybGQhPC9mb250PjwvZGl2Pg==", // "<div style=\"color: red;\"><font style=\"font-size: 12px;\">Hello, World!</font></div>" in base64
							},
						},
					},
				},
			},
			expected: "Hello, World!",
		},
		{
			name: "HTMLエンティティを含むHTMLの場合に正しく処理すること",
			message: &gmail.Message{
				Id: "test-id-4",
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "Subject", Value: "HTMLエンティティテストメール"},
						{Name: "From", Value: "test@example.com"},
						{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0900"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/html",
							Body: &gmail.MessagePartBody{
								Data: "PGRpdj5IZWxsbyAmYW1wOyBXb3JsZCE8L2Rpdj4=", // "<div>Hello &amp; World!</div>" in base64
							},
						},
					},
				},
			},
			expected: "Hello & World!",
		},
		{
			name: "text/plainが優先されること",
			message: &gmail.Message{
				Id: "test-id-5",
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "Subject", Value: "マルチパートテストメール"},
						{Name: "From", Value: "test@example.com"},
						{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0900"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/plain",
							Body: &gmail.MessagePartBody{
								Data: "UGxhaW4gdGV4dCB2ZXJzaW9u", // "Plain text version" in base64
							},
						},
						{
							MimeType: "text/html",
							Body: &gmail.MessagePartBody{
								Data: "PGRpdj5IVE1MIHZlcnNpb248L2Rpdj4=", // "<div>HTML version</div>" in base64
							},
						},
					},
				},
			},
			expected: "Plain text version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.convertToGmailMessage(tt.message)
			assert.Equal(t, tt.expected, result.Body)
			assert.Equal(t, tt.message.Id, result.ID)
		})
	}
}

func TestGmailMessageService_extractBody(t *testing.T) {
	service := &gmailMessageService{}

	tests := []struct {
		name     string
		payload  *gmail.MessagePart
		expected string
	}{
		{
			name: "シンプルなテキスト本文を抽出すること",
			payload: &gmail.MessagePart{
				Body: &gmail.MessagePartBody{
					Data: "SGVsbG8sIFdvcmxkIQ==", // "Hello, World!" in base64
				},
			},
			expected: "Hello, World!",
		},
		{
			name: "マルチパートからtext/plainを抽出すること",
			payload: &gmail.MessagePart{
				Parts: []*gmail.MessagePart{
					{
						MimeType: "text/plain",
						Body: &gmail.MessagePartBody{
							Data: "UGxhaW4gdGV4dA==", // "Plain text" in base64
						},
					},
					{
						MimeType: "text/html",
						Body: &gmail.MessagePartBody{
							Data: "PGRpdj5IVE1MPC9kaXY+", // "<div>HTML</div>" in base64
						},
					},
				},
			},
			expected: "Plain text",
		},
		{
			name: "text/plainがない場合にtext/htmlを抽出すること",
			payload: &gmail.MessagePart{
				Parts: []*gmail.MessagePart{
					{
						MimeType: "text/html",
						Body: &gmail.MessagePartBody{
							Data: "PGRpdj5IVE1MPC9kaXY+", // "<div>HTML</div>" in base64
						},
					},
				},
			},
			expected: "<div>HTML</div>",
		},
		{
			name: "空のペイロードの場合に空文字列を返すこと",
			payload: &gmail.MessagePart{
				Parts: []*gmail.MessagePart{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractBody(tt.payload)
			assert.Equal(t, tt.expected, result)
		})
	}
}
