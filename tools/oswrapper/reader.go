package oswrapper

import (
	"fmt"
	"os"
	"path/filepath"
)

// OsWrapper は両方の機能を持つ具象構造体です
type OsWrapper struct{}

// New は OsWrapper のインスタンスを返します
func New() *OsWrapper {
	return &OsWrapper{}
}

// ReadFile はファイルを読み込み文字列として返します
func (o *OsWrapper) ReadFile(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("パス解決失敗: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("ファイル読み込み失敗: %w", err)
	}

	return string(data), nil
}

// GetEnv は環境変数を取得します
func (o *OsWrapper) GetEnv(key string) string {
	return os.Getenv(key)
}
