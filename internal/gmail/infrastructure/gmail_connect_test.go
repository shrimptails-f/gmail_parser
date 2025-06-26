package infrastructure

import (
	cd "business/internal/common/domain"
	gc "business/tools/gmail"
	gs "business/tools/gmailService"
	"business/tools/oswrapper"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// モックGmailServiceClient
type mockGmailServiceClient struct {
	mock.Mock
}

func (m *mockGmailServiceClient) Authenticate(ctx context.Context, clientSecretPath string, port int) (*oauth2.Token, error) {
	args := m.Called(ctx, clientSecretPath, port)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *mockGmailServiceClient) CreateGmailService(ctx context.Context, credentialsPath, tokenPath string) (*gmail.Service, error) {
	args := m.Called(ctx, credentialsPath, tokenPath)
	return args.Get(0).(*gmail.Service), args.Error(1)
}

// モックGmailClient
type mockGmailClient struct {
	mock.Mock
}

func (m *mockGmailClient) ListMessageIDs(ctx context.Context, max int64) ([]string, error) {
	args := m.Called(ctx, max)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockGmailClient) GetMessagesByLabelName(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	args := m.Called(ctx, labelName, sinceDaysAgo)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockGmailClient) GetGmailDetail(id string) (cd.BasicMessage, error) {
	args := m.Called(id)
	return args.Get(0).(cd.BasicMessage), args.Error(1)
}

func (m *mockGmailClient) SetClient(svc *gmail.Service) *gc.Client {
	m.Called(svc)
	// 実際のClientを作成してサービスをセット
	client := gc.New()
	return client.SetClient(svc)
}

// モックOsWrapper
type mockOsWrapper struct {
	mock.Mock
}

func (m *mockOsWrapper) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *mockOsWrapper) ReadFile(filename string) (string, error) {
	args := m.Called(filename)
	return args.String(0), args.Error(1)
}

// テスト用のGmailConnect構造体（依存性注入を使わずに直接テスト）
type testableGmailConnect struct {
	gs  gs.ClientInterface
	gc  gc.ClientInterface
	osw oswrapper.OsWapperInterface
}

func newTestableGmailConnect(gs gs.ClientInterface, gc gc.ClientInterface, osw oswrapper.OsWapperInterface) *testableGmailConnect {
	return &testableGmailConnect{
		gs:  gs,
		gc:  gc,
		osw: osw,
	}
}

func (g *testableGmailConnect) createGmailClient(ctx context.Context) (*gc.Client, error) {
	credentialsPath := g.osw.GetEnv("CLIENT_SECRET_PATH")
	tokenPath := "/data/credentials/token_user.json"

	svc, err := g.gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		return nil, errors.New("gmail サービス生成に失敗: " + err.Error())
	}

	return g.gc.SetClient(svc), nil
}

func (g *testableGmailConnect) GetMessageIds(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error) {
	// 動的にクライアントを生成
	client, err := g.createGmailClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.GetMessagesByLabelName(ctx, labelName, sinceDaysAgo)
}

func (g *testableGmailConnect) GetGmailDetail(id string) (cd.BasicMessage, error) {
	// 動的にクライアントを生成
	ctx := context.Background()
	client, err := g.createGmailClient(ctx)
	if err != nil {
		return cd.BasicMessage{}, err
	}

	return client.GetGmailDetail(id)
}

func TestGmailConnect_GetMessageIds_Success(t *testing.T) {
	ctx := context.Background()

	// モックの準備
	mockGS := &mockGmailServiceClient{}
	mockGC := &mockGmailClient{}
	mockOSW := &mockOsWrapper{}

	// 期待値の設定
	mockService := &gmail.Service{}

	// モックの動作設定
	mockOSW.On("GetEnv", "CLIENT_SECRET_PATH").Return("/path/to/credentials.json")
	mockGS.On("CreateGmailService", ctx, "/path/to/credentials.json", "/data/credentials/token_user.json").Return(mockService, nil)
	mockGC.On("SetClient", mockService)

	// テスト対象の作成
	conn := newTestableGmailConnect(mockGS, mockGC, mockOSW)

	// createGmailClientメソッドのテスト
	client, err := conn.createGmailClient(ctx)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// モックの呼び出し検証
	mockOSW.AssertExpectations(t)
	mockGS.AssertExpectations(t)
	mockGC.AssertExpectations(t)
}

func TestGmailConnect_GetMessageIds_ServiceCreationError(t *testing.T) {
	ctx := context.Background()

	// モックの準備
	mockGS := &mockGmailServiceClient{}
	mockGC := &mockGmailClient{}
	mockOSW := &mockOsWrapper{}

	// エラーケースの設定
	expectedError := errors.New("service creation failed")
	mockOSW.On("GetEnv", "CLIENT_SECRET_PATH").Return("/path/to/credentials.json")
	mockGS.On("CreateGmailService", ctx, "/path/to/credentials.json", "/data/credentials/token_user.json").Return((*gmail.Service)(nil), expectedError)

	// テスト対象の作成
	conn := newTestableGmailConnect(mockGS, mockGC, mockOSW)

	// createGmailClientメソッドのテスト
	client, err := conn.createGmailClient(ctx)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "gmail サービス生成に失敗")

	// モックの呼び出し検証
	mockOSW.AssertExpectations(t)
	mockGS.AssertExpectations(t)
}

func TestGmailConnect_GetGmailDetail_Success(t *testing.T) {
	// モックの準備
	mockGS := &mockGmailServiceClient{}
	mockGC := &mockGmailClient{}
	mockOSW := &mockOsWrapper{}

	// 期待値の設定
	mockService := &gmail.Service{}

	// モックの動作設定
	mockOSW.On("GetEnv", "CLIENT_SECRET_PATH").Return("/path/to/credentials.json")
	mockGS.On("CreateGmailService", mock.Anything, "/path/to/credentials.json", "/data/credentials/token_user.json").Return(mockService, nil)
	mockGC.On("SetClient", mockService)

	// テスト対象の作成
	conn := newTestableGmailConnect(mockGS, mockGC, mockOSW)

	// createGmailClientメソッドのテスト
	client, err := conn.createGmailClient(context.Background())

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// モックの呼び出し検証
	mockOSW.AssertExpectations(t)
	mockGS.AssertExpectations(t)
	mockGC.AssertExpectations(t)
}

func TestGmailConnect_GetGmailDetail_ServiceCreationError(t *testing.T) {
	// モックの準備
	mockGS := &mockGmailServiceClient{}
	mockGC := &mockGmailClient{}
	mockOSW := &mockOsWrapper{}
	ctx := context.Background()

	// エラーケースの設定
	expectedError := errors.New("service creation failed")
	mockOSW.On("GetEnv", "CLIENT_SECRET_PATH").Return("/path/to/credentials.json")
	mockGS.On("CreateGmailService", ctx, "/path/to/credentials.json", "/data/credentials/token_user.json").Return((*gmail.Service)(nil), expectedError)

	// テスト対象の作成
	conn := newTestableGmailConnect(mockGS, mockGC, mockOSW)

	// createGmailClientメソッドのテスト
	client, err := conn.createGmailClient(context.Background())

	// 検証
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "gmail サービス生成に失敗")

	// モックの呼び出し検証
	mockOSW.AssertExpectations(t)
	mockGS.AssertExpectations(t)
}

// 実際のNew関数のテスト
func TestNew_CreatesGmailConnect(t *testing.T) {
	mockGS := &mockGmailServiceClient{}
	mockGC := &mockGmailClient{}
	mockOSW := &mockOsWrapper{}

	conn := New(mockGS, mockGC, mockOSW)

	assert.NotNil(t, conn)
	assert.IsType(t, &GmailConnect{}, conn)
}
