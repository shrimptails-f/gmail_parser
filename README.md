![Version](https://img.shields.io/badge/Version-1.0.0-green)
# プロジェクトの概要説明
## 営業メールって、正直つらくないですか？
毎日のように届く営業メール<br>
・数が多い<br>
・検索がしづらい<br>
・欲しい情報が載っていないことも多い<br>
そんな不便さを何とかしたくて、IT業界のエンジニア向け営業メールをAIで字句解析して、DBに保存して検索できるアプリを作ってみました。
## 集計イメージ
![image](https://github.com/user-attachments/assets/fee1cd7c-b4c8-428c-806e-2dbfb4eb51a5)
# 環境構築手順
## [VsCode使用者向け](./docs/VsCodeDevContainer.md)
## メール収集方法(環境構築後)
以下のコマンドでGメールを取得→AIで解析→DB保存ができます。
```bash
task gmail-messages-by-label -- 0
```
## 取得結果を表示する
DBに保存したデータの表示方法は[こちら](./docs/query.md) を参照してください。
# 開発者向け情報
## 言語
* Go1.24
## DB
* MySQL8.0
## 環境構築
* Docker
* DevContainer
## ディレクトリ構成の方針
* クリーンアーキテクチャ
## テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./docs/migration.md) を参照してください。
## テスト実施方法
以下のコマンドでカバレッジレポート出力できます。
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```