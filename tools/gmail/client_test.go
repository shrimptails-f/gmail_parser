//go:build integration
// +build integration

package gmail

import (
	"business/tools/gmailService"
	"business/tools/oswrapper"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Gメール認証をテストします。
func TestGmailAuthenticate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	osw := oswrapper.New()
	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")

	strPort := osw.GetEnv("GMAIL_PORT")
	port, err := strconv.Atoi(strPort)
	assert.NoError(t, err)

	gs := gmailService.NewClient()
	_, err = gs.Authenticate(ctx, credentialsPath, port)
	assert.NoError(t, err)
}

// Gメールが取得できるかテストを行います。
func TestRealGmailConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	osw := oswrapper.New()
	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")
	gs := gmailService.NewClient()

	tokenPath := "/data/credentials/token_user.json"
	svc, err := gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		t.Fatalf("gメールAPIクライアント作成失敗: %v", err)
	}

	cilent := NewClient(svc)
	ids, err := cilent.ListMessageIDs(ctx, 5)
	if err != nil {
		t.Fatalf("メールID取得失敗: %v", err)
	}

	if len(ids) == 0 {
		t.Log("メールは見つかりませんでした")
	} else {
		t.Logf("取得したメールID: %v", ids)
		for _, id := range ids {
			fmt.Println("Mail ID:", id)
		}
	}
}

func TestGetGmailByLabel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	osw := oswrapper.New()
	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")
	gs := gmailService.NewClient()

	tokenPath := "/data/credentials/token_user.json"
	svc, err := gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		t.Fatalf("gメールAPIクライアント作成失敗: %v", err)
	}

	client := NewClient(svc)
	messages, err := client.GetMessagesByLabelName(ctx, "営業/テスト")
	if err != nil {
		t.Fatalf("ラベル取得失敗: %v\n", err)
	}
	if len(messages) == 0 {
		t.Log("メールは見つかりませんでした \n")
	} else {
		for _, message := range messages {
			t.Log("\n 取得したメールメッセージ内容:")
			fmt.Printf("メッセージID:%v \n", message.ID)
			fmt.Printf("メッセージ件名:%v \n", message.Subject)
			fmt.Printf("メッセージ送信者:%v \n", message.From)
			fmt.Printf("メッセージ受診:%v \n", message.Date)
			// fmt.Printf("メッセージ本文:%v \n", message.Body)
		}
	}
}
