package gmail

import (
	cd "business/internal/common/domain"
	"context"
)

type ClientInterface interface {
	ListMessageIDs(ctx context.Context, max int64) ([]string, error)
	GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error)
}
