![Version](https://img.shields.io/badge/Version-1.0.0-green)
# プロジェクトの概要説明
### 営業メールって、正直つらくないですか？
毎日のように届く営業メール、
・数が多い
・検索がしづらい
・欲しい情報が載っていないことも多い

そんな不便さを何とかしたくて、IT業界のエンジニア向け営業メールをAIで字句解析して、DBに保存して検索できるアプリを作ってみました。
## 言語
* Go1.24
## DB
* MySQL8.0
## 環境構築
* Docker
* DevContainer
## ディレクトリ構成の方針
* クリーンアーキテクチャ
## 環境構築手順
### [VsCode使用者向け](./docs/VsCodeDevContainer.md) 
## メール収集方法(環境構築後)
以下のコマンドでGメールを取得→AIで解析→DB保存ができます。
```bash
task gmail-messages-by-label -- 0
```
## テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./docs/migration.md) を参照してください。
## テスト実施方法
以下のコマンドでカバレッジレポート出力できます。
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```
