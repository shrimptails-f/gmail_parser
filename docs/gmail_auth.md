# Gmail認証機能

このドキュメントでは、Gmail認証機能の設定と使用方法について説明します。

## 概要

Gmail認証機能は、Google OAuth2を使用してGmailアカウントでの認証を行うコマンドラインアプリケーションです。クリーンアーキテクチャに基づいて実装されており、以下の層に分かれています：

- **ドメイン層**: 認証に関するビジネスルールとドメインモデル
- **アプリケーション層**: 認証のユースケースとインターフェース
- **インフラストラクチャ層**: Google OAuth2 API、データベース、JWTの実装
- **プレゼンテーション層**: コマンドラインインターフェース

## 前提条件

1. Go 1.24以上
2. MySQL データベース
3. Google Cloud Console でのOAuth2設定

## Google OAuth2設定

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. 新しいプロジェクトを作成するか、既存のプロジェクトを選択
3. 「APIとサービス」→「認証情報」に移動
4. 「認証情報を作成」→「OAuth 2.0 クライアントID」を選択
5. アプリケーションの種類として「ウェブアプリケーション」を選択
6. 承認済みのリダイレクトURIに `http://localhost:8080/auth/google/callback` を追加
7. クライアントIDとクライアントシークレットをメモ

## 環境設定

1. `.env.sample`を`.env`にコピー：
```bash
cp .env.sample .env
```

2. `.env`ファイルを編集して必要な値を設定：
```bash
# データベース設定
MYSQL_USER=your_mysql_user
MYSQL_PASSWORD=your_mysql_password
MYSQL_DATABASE=your_database_name
DB_HOST=localhost
DB_PORT=3306

# Google OAuth2設定
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret

# JWT設定
JWT_SECRET_KEY=your_jwt_secret_key_here

# アプリケーション設定
APP_NAME=Gmail Auth App
```

3. 環境変数を読み込み：
```bash
source .env
```

## データベースマイグレーション

ユーザーテーブルを作成するためにマイグレーションを実行：

```bash
go run tools/migrations/main.go
```

## 使用方法

### 1. 認証URLの生成

Google認証URLを生成します：

```bash
cd cmd/gmail_auth
go run main.go auth-url
```

出力例：
```
Google認証URL:
https://accounts.google.com/o/oauth2/auth?client_id=...&redirect_uri=...

ブラウザでこのURLにアクセスして認証を完了してください。
認証後、リダイレクトURLのcodeパラメータを使用して 'auth-code' コマンドを実行してください。
```

### 2. 認証の実行

ブラウザで認証URLにアクセスし、Googleアカウントでログインします。認証後、リダイレクトURLから`code`パラメータを取得し、以下のコマンドを実行：

```bash
go run main.go auth-code <認証コード>
```

成功例：
```
認証成功!
ユーザーID: 1
メールアドレス: user@gmail.com
名前: 山田 太郎
新規ユーザー: true
JWTトークン: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. JWTトークンの検証

生成されたJWTトークンを検証します：

```bash
go run main.go validate-token <JWTトークン>
```

成功例：
```
トークン検証成功!
ユーザーID: 1
```

## アーキテクチャ

### ディレクトリ構造

```
internal/auth/
├── domain/          # ドメイン層
│   └── auth.go      # ドメインモデルとビジネスルール
├── application/     # アプリケーション層
│   ├── interface.go # インターフェース定義
│   └── usecase.go   # ユースケース実装
├── infrastructure/ # インフラストラクチャ層
│   ├── repository.go    # データベースリポジトリ
│   ├── google_oauth.go  # Google OAuth2サービス
│   └── jwt_service.go   # JWTサービス
└── di/             # 依存性注入
    └── wire.go     # 依存関係の設定
```

### 主要コンポーネント

#### ドメイン層
- `User`: ユーザーのドメインモデル
- `GoogleAuthRequest/Response`: Google認証のリクエスト/レスポンス
- `GoogleUserInfo`: Googleから取得するユーザー情報
- `AuthResult`: 認証結果

#### アプリケーション層
- `AuthUseCase`: 認証のユースケース
- `AuthRepository`: ユーザーデータのリポジトリインターフェース
- `GoogleOAuthService`: Google OAuth2サービスインターフェース
- `JWTService`: JWTトークン管理インターフェース

#### インフラストラクチャ層
- `authRepository`: データベースアクセスの実装
- `googleOAuthService`: Google OAuth2 APIの実装
- `jwtService`: JWTトークン生成・検証の実装

## エラーハンドリング

アプリケーションは以下のエラーを適切に処理します：

- `ErrInvalidAuthCode`: 無効な認証コード
- `ErrUserNotFound`: ユーザーが見つからない
- `ErrInvalidToken`: 無効なJWTトークン
- `ErrTokenExpired`: JWTトークンの有効期限切れ
- `ErrEmailNotVerified`: メールアドレスが認証されていない
- `ErrInvalidGoogleUser`: 無効なGoogleユーザー情報

## セキュリティ考慮事項

1. **JWT秘密鍵**: `JWT_SECRET_KEY`は十分に複雑で予測困難な値を使用
2. **OAuth2設定**: Google Cloud ConsoleでリダイレクトURIを適切に設定
3. **環境変数**: 本番環境では`.env`ファイルをバージョン管理に含めない
4. **HTTPS**: 本番環境ではHTTPSを使用
5. **State検証**: 実際の実装ではCSRF攻撃を防ぐためstateパラメータを適切に検証

## トラブルシューティング

### よくある問題

1. **環境変数が設定されていない**
   ```
   GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set
   ```
   → `.env`ファイルを確認し、`source .env`を実行

2. **データベース接続エラー**
   ```
   main - mysql.New: dial tcp: connect: connection refused
   ```
   → MySQLサーバーが起動していることを確認

3. **無効な認証コード**
   ```
   認証に失敗しました: google oauth - ExchangeCode: oauth2: cannot fetch token
   ```
   → 認証コードが正しいか、有効期限内かを確認

4. **JWTトークンエラー**
   ```
   JWT_SECRET_KEY environment variable is not set
   ```
   → `JWT_SECRET_KEY`環境変数を設定

## 開発・テスト

### テストの実行

```bash
# 単体テスト
go test ./internal/auth/...

# 統合テスト
go test -tags=integration ./internal/auth/...
```

### 依存関係の追加

新しい依存関係を追加する場合：

```bash
go mod tidy
```

## 今後の拡張

- Web APIエンドポイントの追加
- リフレッシュトークンの実装
- ユーザー情報の更新機能
- ログアウト機能
- 複数の認証プロバイダー対応
