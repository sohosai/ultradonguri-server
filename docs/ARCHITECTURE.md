# アーキテクチャ設計書

## 概要

本プロジェクトは **Clean Architecture** の思想を取り入れたレイヤードアーキテクチャで構成されています。

- ドメインロジックは `domain` 層に集約し、外部技術（OBS、HTTP、WebSocket、ファイル IO）への依存を排除しています。
- 依存関係は **内側（domain）→ 外側（infrastructure/presentation）** の一方向に保たれています。

```
┌─────────────────────────────────────────────┐
│  Presentation (HTTP Handlers / Requests /   │
│  Responses)                                   │
├─────────────────────────────────────────────┤
│  Infrastructure (OBS Client / WebSocket /   │
│  File IO)                                     │
├─────────────────────────────────────────────┤
│  Domain (Entities / Repository Interfaces)  │
└─────────────────────────────────────────────┘
```

---

## レイヤー構成

### 1. Domain 層 (`src/internal/domain/`)

ビジネスルールとドメインモデルを定義します。

| パッケージ | 責務 |
|-----------|------|
| `entities` | ドメインエンティティ（Performance, Music, TelopMessage, AppError など） |
| `repositories` | リポジトリインターフェース（SceneManager, TelopManager） |

外部技術を一切知らず、**純粋な Go の構造体・インターフェース**のみで記述されています。

#### 主要エンティティ

- **`PerformancePost` / `Performance` / `Music`** — パフォーマンス開始・楽曲情報
- **`ConversionPost` / `NextPerformance`** — 転換中の情報
- **`TelopMessage`** — Viewer へ送信するテロップの統一表現
- **`MuteState` / `CMState` / `DisplayCopyright`** — 各種状態制御

---

### 2. Infrastructure 層 (`src/internal/infrastructure/`)

外部システムとのやり取りを担当します。domain 層で定義されたインターフェースを実装します。

| パッケージ | 責務 |
|-----------|------|
| `scene` | OBS WebSocket API 経由でシーンを取得・切り替え、バックアップ/復旧 |
| `telop` | テロップ状態の保持・更新、バックアップ/復旧 |
| `telop/websocket` | Viewer との WebSocket 接続管理・ブロードキャスト |
| `file` | `events.json` からのパフォーマンス情報読み込み |

#### SceneManager (`scene`)

- OBS の **UUID** を用いてシーンを特定（シーン名は変更可能なため）
- 3 つのシーン（Normal / Muted / CM）を管理
- `SetMute()` は `isForceMutedFlag` を考慮して動作
- 状態（現在のシーン種別・強制ミュートフラグ）を JSON ファイルに永続化

#### TelopManager (`telop`)

- `Either3<PerformancePost, ConversionPost, EmptyTelop>` を用いて排他的に状態を保持
- テロップ変更時に自動で JSON ファイルに永続化
- 起動時にバックアップファイルから復旧可能

#### WebSocketHub (`telop/websocket`)

- `gorilla/websocket` を用いて接続管理
- `telopChannel` を通じた非同期ブロードキャスト
- 接続の追加・削除はミューテックスで排他制御

---

### 3. Presentation 層 (`src/internal/presentation/`)

HTTP リクエストの受付・レスポンス返却を担当します。

| パッケージ | 責務 |
|-----------|------|
| `handlers` | Gin のハンドラ実装（各エンドポイント） |
| `model/requests` | リクエストボディのバインド用構造体 |
| `model/responses` | レスポンス用構造体およびエラー処理ヘルパー |

---

## データフロー

### パフォーマンス開始時

```
[Controller]  POST /performance/start
       │
       ▼
[Handler]  bind request → create domain entity
       │
       ├──► [TelopManager] SetPerformanceTelop()
       │         └── save to telop.json
       │
       ├──► [WebSocketHub] PushTelop() → broadcast to viewers
       │
       └──► [SceneManager] SetNormalScene()
                 └── switch OBS scene + save to scene.json
```

### 楽曲切替時

```
[Controller]  POST /performance/music
       │
       ▼
[Handler]  bind request → create domain entity
       │
       ├──► [TelopManager] SetMusicTelop()
       │         └── update current telop + save to telop.json
       │
       ├──► [WebSocketHub] PushTelop() → broadcast music info
       │
       └──► [SceneManager] SetMute(shouldBeMuted)
                 └── Normal or Muted scene
```

### Viewer 接続時

```
[Viewer]  GET /ws
       │
       ▼
[WebsocketHandlers] upgrade to WebSocket
       │
       ▼
[WebSocketHub] AddConnection()
       │
       └── 以降、PushTelop() により broadcast されるメッセージを受信
```

---

## バックアップとリストア

### リストアフロー（起動時）

1. `main()` で `scene.RestoreSceneManager()` を呼び出し
   - `scene.json` が存在すれば、シーン種別・強制ミュートフラグを復元
   - 存在しない場合は新規作成
2. `main()` で `telop.RestoreTelopManager()` を呼び出し
   - `telop.json` が存在すれば、テロップ種別・内容を復元
   - 存在しない場合は空の Conversion 状態で初期化

### バックアップ保存タイミング

- `SceneManager`: シーン切替時・強制ミュートフラグ変更時に `saveToFile()` を実行
- `TelopManager`: テロップ更新時に `saveToFile()` を `defer` で実行

---

## CORS 設定

`CONTROLLER_ADDRESS` に指定されたオリジンのみを許可します。プレフィックス一致で判定されるため、ポート番号が変わっても許可されます。

```go
AllowOriginFunc: func(origin string) bool {
    for _, prefix := range allowOrigins {
        prefix = strings.TrimSpace(prefix)
        if strings.HasPrefix(origin, prefix) {
            return true
        }
    }
    return false
}
```

---

## 技術的な選定理由

| 技術 | 選定理由 |
|------|----------|
| **goobs** | OBS の WebSocket API v5 にネイティブ対応し、型安全なクライアントを提供 |
| **Either3 / Option** | `samber/mo` の代わりに独自の `utils.Option` を併用。テロップ状態を排他的に表現し、null 安全性を確保 |
| **Air** | 開発時のフィードバックループを短縮 |
| **Swagger (swaggo)** | ハンドラのコメントからドキュメントを自動生成し、メンテナンスコストを削減 |

---

## ファイル構成の補足

- `events.json` — パフォーマンス情報を静的 JSON で管理。DB は使用しない。
- `scene.json` / `telop.json` — ランタイム状態の永続化ファイル。`.gitignore` で Git 管理対象外。
