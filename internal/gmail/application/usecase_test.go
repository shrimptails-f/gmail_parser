package application

import (
	cd "business/internal/common/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGmailConnect はGmailConnectInterfaceのモック実装です
type MockGmailConnect struct {
	mock.Mock
}

func (m *MockGmailConnect) GetMessageIds(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	args := m.Called(ctx, labelName, sinceDaysAgo)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGmailConnect) GetGmailDetail(id string) (cd.BasicMessage, error) {
	args := m.Called(id)
	return args.Get(0).(cd.BasicMessage), args.Error(1)
}

// MockEmailStoreUseCase はEmailStoreUseCaseのモック実装です
type MockEmailStoreUseCase struct {
	mock.Mock
}

func (m *MockEmailStoreUseCase) SaveEmailAnalysisResult(result cd.Email) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockEmailStoreUseCase) GetEmailByGmailIds(gmailIds []string) ([]string, error) {
	args := m.Called(gmailIds)
	return args.Get(0).([]string), args.Error(1)
}

func TestGmailUseCase_GetMessages(t *testing.T) {
	ctx := context.Background()

	// テストデータの準備
	testMessageIds := []string{"msg1", "msg2", "msg3"}
	existingIds := []string{"msg1"} // msg1は既にDBに存在

	testMessage := cd.BasicMessage{
		ID:      "msg2",
		Subject: "Test Subject",
		From:    "from@example.com",
		To:      []string{"to@example.com"},
		Date:    time.Now(),
		Body:    "This is a test body",
	}

	// モックの設定
	mockGmailConnect := &MockGmailConnect{}
	mockEmailStore := &MockEmailStoreUseCase{}

	// GetMessageIds のモック設定
	mockGmailConnect.On("GetMessageIds", ctx, "INBOX", 7).Return(testMessageIds, nil)

	// GetEmailByGmailIds のモック設定（既存のメールIDを返す）
	mockEmailStore.On("GetEmailByGmailIds", testMessageIds).Return(existingIds, nil)

	// GetGmailDetail のモック設定（新しいメールの詳細を返す）
	mockGmailConnect.On("GetGmailDetail", "msg2").Return(testMessage, nil)
	mockGmailConnect.On("GetGmailDetail", "msg3").Return(cd.BasicMessage{
		ID:      "msg3",
		Subject: "Another Test Subject",
		From:    "another@example.com",
		To:      []string{"to@example.com"},
		Date:    time.Now(),
		Body:    "Another test body",
	}, nil)

	// ユースケースの作成
	useCase := New(mockGmailConnect, mockEmailStore)

	// テスト実行
	result, err := useCase.GetMessages(ctx, "INBOX", 7)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result)) // msg2とmsg3の2件が返される（msg1は既存のためスキップ）

	// モックの呼び出し確認
	mockGmailConnect.AssertExpectations(t)
	mockEmailStore.AssertExpectations(t)
}

func TestGmailUseCase_GetMessages_ErrorOnGetMessageIds(t *testing.T) {
	ctx := context.Background()

	// モックの設定
	mockGmailConnect := &MockGmailConnect{}
	mockEmailStore := &MockEmailStoreUseCase{}

	// GetMessageIds でエラーを返すモック設定
	mockGmailConnect.On("GetMessageIds", ctx, "INBOX", 7).Return([]string{}, assert.AnError)

	// ユースケースの作成
	useCase := New(mockGmailConnect, mockEmailStore)

	// テスト実行
	result, err := useCase.GetMessages(ctx, "INBOX", 7)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)

	// モックの呼び出し確認
	mockGmailConnect.AssertExpectations(t)
}

func TestGmailUseCase_GetMessages_ErrorOnGetEmailByGmailIds(t *testing.T) {
	ctx := context.Background()

	// テストデータの準備
	testMessageIds := []string{"msg1", "msg2"}

	// モックの設定
	mockGmailConnect := &MockGmailConnect{}
	mockEmailStore := &MockEmailStoreUseCase{}

	// GetMessageIds のモック設定
	mockGmailConnect.On("GetMessageIds", ctx, "INBOX", 7).Return(testMessageIds, nil)

	// GetEmailByGmailIds でエラーを返すモック設定
	mockEmailStore.On("GetEmailByGmailIds", testMessageIds).Return([]string{}, assert.AnError)

	// ユースケースの作成
	useCase := New(mockGmailConnect, mockEmailStore)

	// テスト実行
	result, err := useCase.GetMessages(ctx, "INBOX", 7)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "GetMessages:")

	// モックの呼び出し確認
	mockGmailConnect.AssertExpectations(t)
	mockEmailStore.AssertExpectations(t)
}
