// Package infrastructure はAI機能のインフラストラクチャ層を提供します。
// このファイルはプロンプトファイルの読み込み機能を実装します。
package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilePromptService はファイルシステムからプロンプトを読み込むサービスです
type FilePromptService struct {
	promptDir string
}

// NewFilePromptService はファイルプロンプトサービスを作成します
func NewFilePromptService(promptDir string) *FilePromptService {
	return &FilePromptService{
		promptDir: promptDir,
	}
}

// LoadPrompt は指定されたファイル名のプロンプトを読み込みます
func (s *FilePromptService) LoadPrompt(filename string) (string, error) {
	// ファイル名の妥当性チェック
	if filename == "" {
		return "", fmt.Errorf("ファイル名が空です")
	}

	// ディレクトリトラバーサル攻撃を防ぐ
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("無効なファイルパス: %s", filename)
	}

	// ファイルパスを構築
	filePath := filepath.Join(s.promptDir, filename)

	// ファイルが指定されたディレクトリ内にあることを確認
	absPromptDir, err := filepath.Abs(s.promptDir)
	if err != nil {
		return "", fmt.Errorf("プロンプトディレクトリの絶対パス取得エラー: %w", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルパスの絶対パス取得エラー: %w", err)
	}

	if !strings.HasPrefix(absFilePath, absPromptDir) {
		return "", fmt.Errorf("無効なファイルパス: %s", filename)
	}

	// ファイルを読み込み
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	return string(content), nil
}

// SavePrompt は指定されたファイル名でプロンプトを保存します
func (s *FilePromptService) SavePrompt(filename, content string) error {
	// ファイル名の妥当性チェック
	if filename == "" {
		return fmt.Errorf("ファイル名が空です")
	}

	// ディレクトリトラバーサル攻撃を防ぐ
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("無効なファイルパス: %s", filename)
	}

	// ファイルパスを構築
	filePath := filepath.Join(s.promptDir, filename)

	// ファイルが指定されたディレクトリ内にあることを確認
	absPromptDir, err := filepath.Abs(s.promptDir)
	if err != nil {
		return fmt.Errorf("プロンプトディレクトリの絶対パス取得エラー: %w", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("ファイルパスの絶対パス取得エラー: %w", err)
	}

	if !strings.HasPrefix(absFilePath, absPromptDir) {
		return fmt.Errorf("無効なファイルパス: %s", filename)
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(s.promptDir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}

	return nil
}

// ListPrompts はプロンプトディレクトリ内のファイル一覧を取得します
func (s *FilePromptService) ListPrompts() ([]string, error) {
	// ディレクトリが存在するかチェック
	if _, err := os.Stat(s.promptDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// ディレクトリ内のファイルを読み取り
	entries, err := os.ReadDir(s.promptDir)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリ読み取りエラー: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
