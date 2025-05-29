# EmailStore パッケージ

メール分析結果をデータベースに保存する機能を提供するパッケージです。

## 概要

このパッケージは、OpenAIによるメール分析結果を構造化されたデータベーステーブルに保存する機能を実装しています。クリーンアーキテクチャに基づいて設計されており、TDD（テスト駆動開発）で実装されています。

## アーキテクチャ

```
internal/emailstore/
├── domain/          # ドメイン層 - ビジネスルールとエンティティ
├── application/     # アプリケーション層 - ユースケース
├── infrastructure/ # インフラストラクチャ層 - データベース実装
└── di/             # 依存性注入
```

## データベーステーブル

### 主要テーブル

- **emails**: 全メール共通の基本情報
- **email_projects**: 案件メール専用の詳細情報
- **entry_timings**: 案件の入場時期（複数）を正規化管理

### キーワード管理テーブル

- **keyword_groups**: 正規化された技術キーワードのマスタ
- **key_words**: キーワードの表記ゆれ管理
- **email_keyword_groups**: メールとキーワードの多対多関連

### ポジション管理テーブル

- **position_groups**: 正規化されたポジション名のマスタ
- **position_words**: ポジションの表記ゆれ管理
- **email_position_groups**: メールとポジションの多対多関連

### 業務種別管理テーブル

- **work_type_groups**: 正規化された業務種別マスタ
- **work_type_words**: 業務表記ゆれ管理
- **email_work_type_groups**: メールと業務種別の多対多関連

## 使用方法

### 基本的な使用例

```go
import (
    "business/internal/emailstore/di"
    "business/tools/mysql"
)

// データベース接続
db, err := mysql.New()
if err != nil {
    log.Fatal(err)
}

// 依存性注入でユースケースを取得
emailStoreUseCase := di.ProvideEmailStoreDependencies(db.DB)

// メール分析結果を保存
err = emailStoreUseCase.SaveEmailAnalysisResult(ctx, analysisResult)
if err != nil {
    log.Printf("メール保存エラー: %v", err)
}
```

## 機能

### 保存機能

- **案件メール**: 詳細な案件情報（単価、勤務地、技術要素など）を関連テーブルに保存
- **営業メール**: 基本情報のみを保存
- **キーワード正規化**: 技術キーワードを自動的に正規化して保存
- **重複チェック**: 同一メールIDの重複保存を防止

### データ構造

- **トランザクション管理**: 関連データの整合性を保証
- **正規化設計**: キーワードやポジションの表記ゆれに対応
- **一覧画面対応**: カンマ区切り文字列での高速検索をサポート

## テスト

### 単体テスト

```bash
go test ./internal/emailstore/application/... -v
```

### 統合テスト

```bash
go test -tags=integration ./internal/emailstore/... -v
```

### テスト構成

- **ユースケーステスト**: モックを使用したビジネスロジックのテスト
- **リポジトリテスト**: 実際のデータベースを使用した統合テスト
- **統合テスト**: 全体の動作を確認するエンドツーエンドテスト

## エラーハンドリング

- `ErrEmailNotFound`: メールが見つからない場合
- `ErrEmailAlreadyExists`: メールが既に存在する場合
- `ErrInvalidEmailData`: 無効なメールデータの場合

## 依存関係

- **GORM**: ORM（Object-Relational Mapping）
- **MySQL**: データベース
- **testify**: テストフレームワーク

## 設計原則

1. **クリーンアーキテクチャ**: 依存関係の方向を制御
2. **TDD**: テストファーストでの開発
3. **単一責任の原則**: 各層が明確な責任を持つ
4. **依存性注入**: テスト容易性と柔軟性を確保
5. **エラーハンドリング**: 適切なエラー伝播と文脈の追加

## 今後の拡張

- ポジション情報の自動抽出・保存
- 業務種別の自動分類・保存
- 検索機能の追加
- キャッシュ機能の実装
