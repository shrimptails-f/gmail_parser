// Package application は認証機能のアプリケーション層を提供します。
// このファイルはGmail認証ユースケースの実装を定義します。
package application

import (
	"business/internal/gmail/domain"
	r "business/internal/gmail/infrastructure"
	"context"
)

// GmailAuthUseCaseImpl はGmailAuthUseCaseの実装です
type GmailAuthUseCaseImpl struct {
	gmailAuthService r.GmailAuthService
}

// NewGmailAuthUseCase はGmailAuthUseCaseの新しいインスタンスを作成します
func NewGmailAuthUseCase(gmailAuthService r.GmailAuthService) GmailAuthUseCase {
	return &GmailAuthUseCaseImpl{
		gmailAuthService: gmailAuthService,
	}
}

// AuthenticateGmail はGmail認証を実行します
func (u *GmailAuthUseCaseImpl) AuthenticateGmail(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error) {
	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 既存の認証情報を確認
	existingCredential, err := u.gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err == nil && existingCredential != nil && existingCredential.IsValid() {
		// 既存の有効な認証情報がある場合はそれを使用
		return &domain.GmailAuthResult{
			Credential:      *existingCredential,
			ApplicationName: config.ApplicationName,
			IsNewAuth:       false,
		}, nil
	}

	// 新しい認証を実行
	result, err := u.gmailAuthService.Authenticate(ctx, config)
	if err != nil {
		return nil, err
	}

	// 認証情報を保存
	// if err := u.gmailAuthService.SaveCredentials(config.CredentialsFolder, config.UserID, result.Credential); err != nil {
	// 保存に失敗しても認証結果は返す（警告レベル）
	// 実際の実装ではログに記録する
	// }

	return result, nil
}

// CreateGmailService はGmail APIサービスを作成します
func (u *GmailAuthUseCaseImpl) CreateGmailService(ctx context.Context, config domain.GmailAuthConfig) (interface{}, error) {
	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 認証を実行
	authResult, err := u.AuthenticateGmail(ctx, config)
	if err != nil {
		return nil, err
	}

	// Gmail APIサービスを作成
	service, err := u.gmailAuthService.CreateGmailService(ctx, authResult.Credential, config.ApplicationName)
	if err != nil {
		return nil, err
	}

	return service, nil
}
