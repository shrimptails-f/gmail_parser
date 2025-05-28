# 概要
Vscode使用者向けのDevContainerを使用する際の手順です。  

# 環境構築手順
1. ソースをクローン
```bash
git clone https://github.com/shirimptail.f/business.git
```
2. .envをコピー
```bash
cp .devcontainer/.env.sample .devcontainer/.env
```
3. VsCodeでプロジェクトを開く。
4. Reopen in Containerを押下。
![image](https://github.com/shrimpTails/go_clean_architecture/assets/162465105/48c41f20-ebf3-4a92-ae2e-609a060f7a08)
もし表示されない場合は Ctrl Shift P→Reopen in containerと入力して実行でもおｋ
5. サーバー起動
```
task air
```
このような表示が出たら環境構築完了。  
![image](https://github.com/shirimptail.f/business/assets/73060776/54a74657-e32a-42ab-9c1d-64fea294b58d)
# テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./migration.md) を参照してください。

# デバッグ手順
VsCode向けのデバッグ手順は準備中です。