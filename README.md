![Version](https://img.shields.io/badge/Version-1.0.0-green)
# プロジェクトの概要説明
Go言語の技術的な検証や動作確認用として作成しました。
## 言語
* Go1.22
## DB
* MySQL8.0
## 環境構築
* Docker
* DevContainer
## 開発支援ツール
* Air(ホットリロード)
* delve(デバッガ)
## サポートされているIDE
* VsCode
* GoLand
## DI
* uber-go/dig
## ディレクトリ構成の方針
* クリーンアーキテクチャ
# 環境構築手順
 - VsCode使用者向け  
[DevContainerを使う手順](./docs/VsCodeDevContainer.md)  
# テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./docs/migration.md) を参照してください。
# デバッグ手順
プロジェクトのデバッグ方法については、[デバッグ手順](./docs/debug.md) を参照してください。このドキュメントでは、delveを使用した効率的なデバッグプロセスを紹介しています。

# Gmail認証機能の実行手順

## 前提条件
1. Google Cloud Consoleで OAuth2 認証情報を作成し、`client_secret.json` ファイルをプロジェクトルートに配置してください
2. OAuth2 認証情報の設定で、リダイレクトURIに `http://localhost:5555/Callback` を追加してください

## アプリケーションのビルド
```bash
go build -o gmail_auth ./cmd/gmail_auth
```

## 使用方法

### 基本的な実行方法（client_secret.jsonがプロジェクトルートにある場合）
```bash
# Gmail認証の実行
./gmail_auth gmail-auth

# Gmail APIサービスのテスト
./gmail_auth gmail-service

# ヘルプの表示
./gmail_auth --help
```

### 環境変数を指定した実行方法
```bash
# 環境変数でclient_secret.jsonのパスを指定して実行
CLIENT_SECRET_PATH=/data/client_secret.json ./gmail_auth gmail-auth
CLIENT_SECRET_PATH=/data/client_secret.json ./gmail_auth gmail-service
```

### 実行コマンド例
```bash
# 1. アプリケーションをビルド
go build -o gmail_auth ./cmd/gmail_auth

# 2. Gmail認証を実行（client_secret.jsonがプロジェクトルートにある場合）
./gmail_auth gmail-auth

# 3. または、環境変数でパスを指定して実行
CLIENT_SECRET_PATH=/data/client_secret.json ./gmail_auth gmail-auth
```

## 環境変数
- `CLIENT_SECRET_PATH`: client_secret.jsonファイルのパスを指定（デフォルト: ./client_secret.json）

## 注意事項
- 初回実行時はブラウザでGoogle認証が必要です
- 認証情報は `credentials/` フォルダに保存されます
- Gmail API の読み取り専用スコープを使用します

# テスト実施方法
プロジェクトのテストを実行するには、以下の手順に従ってください。

## 全テストの実行
すべてのパッケージのテストを実行するには、以下のコマンドを使用します：

```bash
go test ./...
```

## エラーがあるテストをスキップする
テスト実行中にエラーが発生しても続行するには、`-failfast=false`フラグを使用します：

```bash
go test -failfast=false ./...
```

## テストカバレッジの確認
テストカバレッジを計測し、HTMLレポートとして出力するには、以下のコマンドを実行します：

```bash
# カバレッジプロファイルを生成
go test -coverprofile=coverage.out ./...

# HTMLレポートを生成
go tool cover -html=coverage.out -o coverage.html
```

生成された`coverage.html`ファイルをブラウザで開くと、コードカバレッジの詳細を視覚的に確認できます。カバレッジレポートでは、テストでカバーされているコード（緑色）とカバーされていないコード（赤色）が表示されます。
