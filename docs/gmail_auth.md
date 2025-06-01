# Gmail認証機能
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