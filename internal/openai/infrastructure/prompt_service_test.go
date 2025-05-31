// Package infrastructure はプロンプトファイル読み込みサービスのテストを提供します。
package infrastructure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilePromptService_LoadPrompt_正常系_プロンプトファイルを読み込むこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	promptContent := "以下のメールを字句解析してください：\n詳細な分析を行い、JSON形式で結果を返してください。"
	promptFile := filepath.Join(tempDir, "text_analysis_prompt.txt")

	err := os.WriteFile(promptFile, []byte(promptContent), 0644)
	assert.NoError(t, err)

	// Act
	result, err := service.LoadPrompt("text_analysis_prompt.txt")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, promptContent, result)
}

func TestFilePromptService_LoadPrompt_異常系_存在しないファイルでエラーを返すこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	// Act
	result, err := service.LoadPrompt("nonexistent.txt")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "ファイル読み込みエラー")
}

func TestFilePromptService_LoadPrompt_異常系_空のファイル名でエラーを返すこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	// Act
	result, err := service.LoadPrompt("")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "ファイル名が空です")
}

func TestFilePromptService_LoadPrompt_異常系_ディレクトリトラバーサル攻撃を防ぐこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	// Act
	result, err := service.LoadPrompt("../../../etc/passwd")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "無効なファイルパス")
}

func TestFilePromptService_LoadPrompt_正常系_UTF8エンコードのファイルを読み込むこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	// 日本語を含むプロンプト
	promptContent := "以下の日本語メールを解析してください：\n感情分析、キーワード抽出、要約を行ってください。"
	promptFile := filepath.Join(tempDir, "japanese_prompt.txt")

	err := os.WriteFile(promptFile, []byte(promptContent), 0644)
	assert.NoError(t, err)

	// Act
	result, err := service.LoadPrompt("japanese_prompt.txt")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, promptContent, result)
	assert.Contains(t, result, "日本語")
}

func TestFilePromptService_LoadPrompt_正常系_大きなファイルを読み込むこと(t *testing.T) {
	t.Parallel()
	// Arrange
	tempDir := t.TempDir()
	service := NewFilePromptService(tempDir)

	// 大きなプロンプトファイルを作成（10KB程度）
	var largeContent string
	for i := 0; i < 1000; i++ {
		largeContent += "これは大きなプロンプトファイルのテストです。\n"
	}

	promptFile := filepath.Join(tempDir, "large_prompt.txt")
	err := os.WriteFile(promptFile, []byte(largeContent), 0644)
	assert.NoError(t, err)

	// Act
	result, err := service.LoadPrompt("large_prompt.txt")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, largeContent, result)
	assert.Greater(t, len(result), 10000) // 10KB以上
}
