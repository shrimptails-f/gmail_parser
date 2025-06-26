// Package infrastructure はGメールとの疎通部分を実装します。
package infrastructure

import (
	cd "business/internal/common/domain"
	gc "business/tools/gmail"
	gs "business/tools/gmailService"
	"business/tools/oswrapper"
	"context"
	"fmt"
)

// GmailConnect Gメール接続処理を持つ構造体です。
type GmailConnect struct {
	gs  gs.ClientInterface
	gc  gc.ClientInterface
	osw oswrapper.OsWapperInterface
}

func New(gs gs.ClientInterface, gc gc.ClientInterface, osw oswrapper.OsWapperInterface) *GmailConnect {
	return &GmailConnect{
		gs:  gs,
		gc:  gc,
		osw: osw,
	}
}

func (g *GmailConnect) createGmailClient(ctx context.Context) (*gc.Client, error) {
	credentialsPath := g.osw.GetEnv("CLIENT_SECRET_PATH")
	tokenPath := "/data/credentials/token_user.json"

	svc, err := g.gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("gmail サービス生成に失敗: %w", err)
	}

	return g.gc.SetClient(svc), nil
}

func (g *GmailConnect) GetMessageIds(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	// 動的にクライアントを生成
	client, err := g.createGmailClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.GetMessagesByLabelName(ctx, labelName, sinceDaysAgo)
}

func (g *GmailConnect) GetGmailDetail(id string) (cd.BasicMessage, error) {
	// 動的にクライアントを生成
	ctx := context.Background()
	client, err := g.createGmailClient(ctx)
	if err != nil {
		return cd.BasicMessage{}, err
	}

	return client.GetGmailDetail(id)
}
