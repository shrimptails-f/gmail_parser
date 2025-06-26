package gmailService

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
}

// GメールのAPIコールをする前の処理をまとめた構造体です
// 認可はURLを手動で開く必要があります。
func New() *Client {
	return &Client{}
}

func (c *Client) Authenticate(ctx context.Context, clientSecretPath string, port int) (*oauth2.Token, error) {
	scopes := []string{"https://www.googleapis.com/auth/gmail.readonly"}

	b, err := os.ReadFile(clientSecretPath)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("ブラウザでこのURLにアクセスして認証してください:\n%v\n", authURL)

	codeCh := make(chan string)
	errCh := make(chan error)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("認証コードが取得できませんでした")
			return
		}
		fmt.Fprintln(w, "認証が完了しました。このウィンドウを閉じてください。")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("認証がタイムアウトしました")
	}

	_ = server.Shutdown(ctx)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// 保存処理
	if err := saveTokenToFile(token); err != nil {
		return nil, err
	}

	return token, nil
}

func (c *Client) CreateGmailService(ctx context.Context, credentialsPath, tokenPath string) (*gmail.Service, error) {
	credBytes, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("クレデンシャル読み込み失敗: %w", err)
	}

	config, err := google.ConfigFromJSON(credBytes, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("OAuth2構成失敗: %w", err)
	}

	token, err := tokenFromFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("トークン読み込み失敗: %w", err)
	}

	svc, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("gmailサービス初期化失敗: %w", err)
	}

	return svc, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveTokenToFile(token *oauth2.Token) error {
	folder := "/data/credentials"
	if err := os.MkdirAll(folder, 0700); err != nil {
		return fmt.Errorf("フォルダ作成失敗: %w", err)
	}

	path := filepath.Join(folder, "token_user.json")

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("トークンのJSON化失敗: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("ファイル保存失敗: %w", err)
	}

	return nil
}
