package gmail

import (
	cd "business/internal/common/domain"
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
)

type Client struct {
	svc *gmail.Service
}

func New() *Client {
	return &Client{}
}

func (c *Client) SetClient(svc *gmail.Service) *Client {
	return &Client{
		svc: svc,
	}
}

func (c *Client) ListMessageIDs(ctx context.Context, max int64) ([]string, error) {
	user := "me"
	resp, err := c.svc.Users.Messages.List(user).MaxResults(max).Do()
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, m := range resp.Messages {
		ids = append(ids, m.Id)
	}
	return ids, nil
}
func (c *Client) GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	user := "me"

	// ラベルID取得
	labelResp, err := c.svc.Users.Labels.List(user).Do()
	if err != nil {
		return nil, fmt.Errorf("ラベル取得に失敗しました。: %v", err)
	}
	var labelID string
	for _, label := range labelResp.Labels {
		if label.Name == labelName {
			labelID = label.Id
			break
		}
	}
	if labelID == "" {
		return nil, fmt.Errorf("ラベル '%s' が見つかりませんでした", labelName)
	}

	// 検索条件
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if sinceDaysAgo != 0 {
		start = start.AddDate(0, 0, sinceDaysAgo)
	}
	query := fmt.Sprintf("after:%d", start.Unix())

	// ページングしながら取得
	var messageIds []string
	pageToken := ""

	for {
		req := c.svc.Users.Messages.List(user).
			LabelIds(labelID).
			Q(query).
			MaxResults(100)
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		resp, err := req.Do()
		if err != nil {
			return nil, err
		}

		for _, m := range resp.Messages {
			messageIds = append(messageIds, m.Id)
		}
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken

	}

	return messageIds, nil
}

func (c *Client) GetGmailDetail(id string) (cd.BasicMessage, error) {
	user := "me"
	full, err := c.svc.Users.Messages.Get(user, id).Format("full").Do()
	if err != nil {
		return cd.BasicMessage{}, fmt.Errorf("gメール取得処理でエラーが発生しました。 %v", err)
	}

	msg := cd.BasicMessage{
		ID:      full.Id,
		Subject: getHeader(full.Payload.Headers, "Subject"),
		From:    getHeader(full.Payload.Headers, "From"),
		To:      parseHeaderMulti(getHeader(full.Payload.Headers, "To")),
		Date:    parseDate(getHeader(full.Payload.Headers, "Date")),
		Body:    stripHTMLTags(extractBody(full.Payload)), // HTMLタグを削除する。
	}
	return msg, nil
}

func getHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, h := range headers {
		if h.Name == name {
			return h.Value
		}
	}
	return ""
}

func parseHeaderMulti(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

func parseDate(raw string) time.Time {
	t, err := mail.ParseDate(raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

func extractBody(payload *gmail.MessagePart) string {
	if (payload.MimeType == "text/plain" || payload.MimeType == "text/html") &&
		payload.Body != nil &&
		payload.Body.Data != "" {

		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)

		if err == nil {
			return string(decoded)
		}
	}
	for _, part := range payload.Parts {
		if body := extractBody(part); body != "" {
			return body
		}
	}
	return ""
}
