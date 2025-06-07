package gmailService

import (
	"business/tools/oswrapper"
	"context"
	"testing"
	"time"
)

func TestRealGmailConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tokenPath := "/data/credentials/token_user.json"
	osw := oswrapper.New()
	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")
	gs := NewClient()
	_, err := gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		t.Fatalf("gメールAPIクライアント作成失敗: %v", err)
	}
}
