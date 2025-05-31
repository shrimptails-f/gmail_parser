// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルはGmail認証サービスの実装を定義します。
package infrastructure

import (
	"business/internal/gmail/application"
	"business/internal/gmail/domain"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// gmailAuthService はGmail認証サービスの実装です
type gmailAuthService struct{}

// NewGmailAuthService は新しいGmail認証サービスを作成します
func NewGmailAuthService() application.GmailAuthService {
	return &gmailAuthService{}
}

// Authenticate はGmail認証を実行します
func (s *gmailAuthService) Authenticate(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error) {
	// client-secret.jsonファイルを読み込み
	clientSecretData, err := ioutil.ReadFile(config.ClientSecretPath)
	if err != nil {
		return nil, fmt.Errorf("client-secret.jsonファイルの読み込みに失敗しました: %w", err)
	}

	// グーグルのOAuth2設定インスタンスを作成
	oauthConfig, err := google.ConfigFromJSON(clientSecretData, config.Scopes...)
	if err != nil {
		return nil, fmt.Errorf("OAuth2設定の作成に失敗しました: %w", err)
	}

	// 認証フローを実行
	token, err := s.getTokenFromWeb(ctx, oauthConfig, config.LocalServerPort)
	if err != nil {
		return nil, fmt.Errorf("認証フローの実行に失敗しました: %w", err)
	}

	// ドメインモデルに変換
	credential := domain.GmailCredential{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresAt:    token.Expiry,
	}

	result := &domain.GmailAuthResult{
		Credential:      credential,
		ApplicationName: config.ApplicationName,
		IsNewAuth:       true,
	}

	return result, nil
}

// CreateGmailService はGmail APIサービスを作成します
func (s *gmailAuthService) CreateGmailService(ctx context.Context, credential domain.GmailCredential, applicationName string) (interface{}, error) {
	// OAuth2トークンを作成
	token := &oauth2.Token{
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		TokenType:    credential.TokenType,
		Expiry:       credential.ExpiresAt,
	}

	// OAuth2設定を作成（最小限の設定）
	config := &oauth2.Config{
		Endpoint: google.Endpoint,
	}

	// HTTPクライアントを作成
	client := config.Client(ctx, token)

	// Gmail APIサービスを作成
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	return service, nil
}

// LoadCredentials は保存された認証情報を読み込みます
func (s *gmailAuthService) LoadCredentials(credentialsFolder, userID string) (*domain.GmailCredential, error) {
	tokenPath := filepath.Join(credentialsFolder, fmt.Sprintf("token_%s.json", userID))

	// ファイルが存在するかチェック
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return nil, domain.ErrClientSecretNotFound
	}

	// ファイルを読み込み
	data, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("認証情報ファイルの読み込みに失敗しました: %w", err)
	}

	// JSONをパース
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("認証情報のパースに失敗しました: %w", err)
	}

	// ドメインモデルに変換
	credential := &domain.GmailCredential{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresAt:    token.Expiry,
	}

	return credential, nil
}

// SaveCredentials は認証情報を保存します
func (s *gmailAuthService) SaveCredentials(credentialsFolder, userID string, credential domain.GmailCredential) error {
	// ディレクトリを作成
	if err := os.MkdirAll(credentialsFolder, 0700); err != nil {
		return fmt.Errorf("認証情報フォルダの作成に失敗しました: %w", err)
	}

	// OAuth2トークンに変換
	token := &oauth2.Token{
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		TokenType:    credential.TokenType,
		Expiry:       credential.ExpiresAt,
	}

	// JSONに変換
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("認証情報のJSONエンコードに失敗しました: %w", err)
	}

	// ファイルに保存
	tokenPath := filepath.Join(credentialsFolder, fmt.Sprintf("token_%s.json", userID))
	if err := ioutil.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("認証情報ファイルの保存に失敗しました: %w", err)
	}

	return nil
}

// getTokenFromWeb はWebブラウザを使用してOAuth2トークンを取得します
func (s *gmailAuthService) getTokenFromWeb(ctx context.Context, config *oauth2.Config, port int) (*oauth2.Token, error) {
	// 認証URLを生成
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("ブラウザでこのURLにアクセスして認証を完了してください:\n%v\n", authURL)

	// ローカルサーバーを起動して認証コードを受信
	codeCh := make(chan string)
	errCh := make(chan error)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("認証コードが取得できませんでした")
			return
		}

		fmt.Fprintf(w, "認証が完了しました。このウィンドウを閉じてください。")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("ローカルサーバーの起動に失敗しました: %w", err)
		}
	}()

	// 認証コードまたはエラーを待機
	var code string
	select {
	case code = <-codeCh:
		// 認証コードを受信
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("認証がタイムアウトしました")
	}

	// サーバーを停止
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("サーバーの停止に失敗しました: %v\n", err)
	}

	// 認証コードをトークンに交換
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("認証コードのトークン交換に失敗しました: %w", err)
	}

	return token, nil
}
