# ultradonguri-server

テロップ送信・OBS シーン制御 API サーバーです。

本サーバーは **OBS Studio** と連携し、WebSocket を通じてクライアント（Viewer）にリアルタイムでテロップ情報を配信します。さらに、パフォーマンス（出演・演奏）の進行に応じて OBS のシーンを自動的に切り替え、ミュート制御や CM モード切替などを行います。

---

## 主な機能

- **テロップ配信** — WebSocket で Viewer にパフォーマンス情報・転換情報・著作権表示をリアルタイム配信
- **OBS シーン制御** — パフォーマンス開始 / 転換 / CM などの状態に応じて OBS のシーンを自動切替
- **ミュート管理** — 楽曲ごとのミュート要件に基づき自動ミュート、および運用者による強制ミュート
- **パフォーマンス情報提供** — `events.json` から出演者・楽曲情報を提供
- **Swagger UI** — インタラクティブな API ドキュメントを `/swagger/index.html` で提供
- **状態の永続化** — シーン状態・テロップ状態をファイルにバックアップし、サーバー再起動時に復旧

---

## 技術スタック

| 項目 | 採用技術 |
|------|----------|
| 言語 | Go 1.24+
| Web フレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| OBS 連携 | [goobs](https://github.com/andreykaipov/goobs) |
| WebSocket | [gorilla/websocket](https://github.com/gorilla/websocket) |
| ドキュメント | [swaggo](https://github.com/swaggo/swag) |
| 開発用ホットリロード | [Air](https://github.com/air-verse/air) |
| コンテナ | Docker / Docker Compose |

---

## ディレクトリ構成

```
.
├── .github/workflows/      # GitHub Actions（GHCR への自動ビルド・プッシュ）
├── src/
│   ├── docs/               # Swagger 自動生成ファイル
│   ├── internal/
│   │   ├── domain/         # ドメイン層（エンティティ・リポジトリインターフェース）
│   │   ├── infrastructure/ # インフラ層（OBS 連携・ファイル IO・WebSocket）
│   │   ├── presentation/   # プレゼンテーション層（HTTP ハンドラ・リクエスト/レスポンスモデル）
│   │   └── utils/          # 汎用ユーティリティ
│   ├── go.mod / go.sum
│   ├── main.go             # エントリーポイント
│   └── .air.toml           # Air 設定
├── compose.yaml            # Docker Compose（開発用）
├── dev.Dockerfile          # 開発用イメージ（Air 付き）
├── prod.Dockerfile         # 本番用イメージ
├── events.json             # パフォーマンス情報サンプル（.gitignore で除外推奨）
├── .sample.env             # 環境変数サンプル
└── README.md
```

---

## セットアップ

### 必要なもの

- [Go](https://go.dev/dl/) 1.24 以上（ローカル実行時）
- [Docker](https://www.docker.com/) + Docker Compose（コンテナ実行時）
- **OBS Studio**（WebSocket サーバー機能を有効化しておく）

### 環境変数

`.env` ファイルをプロジェクトルートに作成してください。`.sample.env` をコピーして編集するのが簡単です。

```bash
cp .sample.env .env
```

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `ADDRESS` | OBS WebSocket サーバーのアドレス | `localhost:4455` |
| `PASSWORD` | OBS WebSocket の接続パスワード | `yourpassword` |
| `NORMAL_SCENE_NAME` | 通常時に使用する OBS シーン名 | `Normal_Scene` |
| `MUTED_SCENE_NAME` | ミュート時に使用する OBS シーン名 | `Muted_Scene` |
| `CM_SCENE_NAME` | CM 時に使用する OBS シーン名 | `CM_Scene` |
| `CONTROLLER_ADDRESS` | CORS を許可するオリジン（カンマ区切り） | `http://localhost,http://127.0.0.1` |
| `SCENE_BACKUP_PATH` | シーン状態のバックアップファイルパス | `scene.json` |
| `TELOP_BACKUP_PATH` | テロップ状態のバックアップファイルパス | `telop.json` |

> **注意**: `SCENE_BACKUP_PATH` と `TELOP_BACKUP_PATH` は `.gitignore` に含まれているため、誤って Git 管理されることはありません。

### ローカル実行

```bash
cd src
go mod download
go run main.go
```

サーバーは `http://0.0.0.0:8080` で起動します。

### Docker で実行（開発用）

```bash
# 開発用（Air によるホットリロード付き）
docker compose up
```

ソースコードを変更すると自動的にビルド・再起動されます。

### Docker で実行（本番用）

```bash
docker build -f prod.Dockerfile -t ultradonguri-server .
docker run -p 8080:8080 --env-file .env ultradonguri-server
```

---

## API 概要

ベースパス: `/`

| メソッド | パス | 概要 |
|----------|------|------|
| `GET` | `/health` | ヘルスチェック |
| `GET` | `/performances` | 全パフォーマンス情報を取得 |
| `POST` | `/force_mute` | 強制ミュートの ON/OFF |
| `POST` | `/performance/start` | パフォーマンス開始（テロップ設定 + Normal シーン） |
| `POST` | `/performance/music` | 楽曲開始（テロップ更新 + 自動ミュート） |
| `POST` | `/conversion/start` | 転換開始（テロップ設定 + シーン切替） |
| `POST` | `/conversion/cm-mode` | CM モードの ON/OFF（転換中のみ有効） |
| `POST` | `/display-copyright` | 著作権表示の ON/OFF |
| `GET` | `/ws` | Viewer 用 WebSocket 接続 |
| `GET` | `/swagger/*any` | Swagger UI（API ドキュメント） |

詳細なリクエスト/レスポンス仕様、および各種スキーマは **Swagger UI** (`http://localhost:8080/swagger/index.html`) または `docs/API.md` を参照してください。

---

## WebSocket 概要

Viewer は `/ws` へ WebSocket 接続することで、リアルタイムにテロップ情報を受信できます。

送信されるメッセージは以下の 5 種類です：

| Type | 説明 |
|------|------|
| `/performance/start` | パフォーマンス開始情報 |
| `/performance/music` | 楽曲情報（ミュート状態含む） |
| `/conversion/start` | 次のパフォーマンス情報（転換中） |
| `/conversion/cm-mode` | CM モードの ON/OF |
| `/display-copyright` | 著作権表示の ON/OFF |

詳細なメッセージ形式は `docs/WEBSOCKET.md` を参照してください。

---

## バックアップと復旧

本サーバーは以下の 2 つの状態をファイルに自動保存し、起動時に復旧します。

- **シーン状態** (`SCENE_BACKUP_PATH`): 現在のシーン種別（Normal / Muted / CM）と強制ミュートフラグ
- **テロップ状態** (`TELOP_BACKUP_PATH`): 現在のテロップ種別（Performance / Conversion / Empty）とその内容

これにより、サーバーが再起動しても直前の運用状態を維持できます。

---

## CI/CD

`.github/workflows/build-and-push.yaml` により、以下のタイミングで Docker イメージが自動ビルド・プッシュされます。

- `main` ブランチへの push → `ghcr.io/sohosai/ultradonguri-server:production`
- `develop` ブランチへの push → `ghcr.io/sohosai/ultradonguri-server:develop`
- 手動実行 (`workflow_dispatch`) → ブランチ名に応じたタグ

---

## ドキュメント一覧

- `docs/API.md` — REST API 詳細仕様
- `docs/ARCHITECTURE.md` — アーキテクチャとデータフロー
- `docs/WEBSOCKET.md` — WebSocket メッセージ仕様
- `src/docs/swagger.yaml` — Swagger 定義（自動生成）

---

## ライセンス

未定
