# 環境構築手順
## openai APIキーを発行
1. https://openai.com/ja-JP/api/ にアクセスしてアカウント登録
2. https://platform.openai.com/api-keys にアクセスしてAPIキーを発行
3. 発行したAPIキーを控えておく

## 注意事項
- OpenAI APIの利用には料金が発生します
- APIキーは安全に管理してください

## Google OAuth2設定
1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. https://console.developers.google.com/apis/library にアクセスし、「Gmail API」で検索→Gmail-apiを有効化する。
4. 「https://console.cloud.google.com/apis/credentials にアクセスし、認証情報を作成→OAuthクライアントIDを選択する。
5. アプリケーションの種類として「ウェブアプリケーション」を選択
6. 承認済みのリダイレクトURIに `http://localhost:5555/Callback` を追加
7. 秘密鍵をJSONでダウンロード あとで使います。

## ソースをクローン
```bash
git clone https://github.com/shrimptails-f/gmail_parser.git
```
## .envをコピー
```bash
cp .devcontainer/.env.sample .devcontainer/.env
```
## プロンプトファイルをコピー
プロンプトファイルの内容を変更することで調整できます
```bash
cp /data/prompts/text_analysis_prompt_sample.txt /data/prompts/text_analysis_prompt.txt
```

### 環境変数設定
`.env`ファイルを編集して必要な値を設定：
```bash
CLIENT_SECRET_PATH=/data/ダウンロードしたファイル名を記載
OPENAI_API_KEY=生成したOpenAiのAPIキーを記載
LABEL=案件メールが振分済のラベル名を記載
GMAIL_PORT=5555
```
## VsCodeでプロジェクトフォルダーを開く
## Reopen in Containerを押下
もし表示されない場合は Ctrl Shift P→Reopen in containerと入力して実行でもおｋ
## テーブル作成
```bash
task migration-create
```
### 認証URLの生成
Google認証URLを生成します：
```bash
task gmail-auth
```
出力例：
```
Google認証URL:
https://accounts.google.com/o/oauth2/auth?client_id=...&redirect_uri=...

ブラウザでこのURLにアクセスして認証を完了してください。
認証後、リダイレクトURLのcodeパラメータを使用して 'auth-code' コマンドを実行してください。
```

# 環境構築完了です！！
お疲れ様でした。