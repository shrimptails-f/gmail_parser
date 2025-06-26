package gmail

import (
	cd "business/internal/common/domain"
	"context"

	"google.golang.org/api/gmail/v1"
)

type ClientInterface interface {
	ListMessageIDs(ctx context.Context, max int64) ([]string, error)
	GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error)
	GetGmailDetail(id string) (cd.BasicMessage, error)
	SetClient(svc *gmail.Service) *Client
}
