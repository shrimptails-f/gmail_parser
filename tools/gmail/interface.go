package gmail

import (
	"context"
)

type ClientInterface interface {
	ListMessageIDs(ctx context.Context, max int64) ([]string, error)
}
