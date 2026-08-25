# API 仕様書

ベース URL: `http://{host}:8080`

CORS は環境変数 `CONTROLLER_ADDRESS` で指定されたオリジンのみを許可しています。

すべての POST/PUT リクエストは `Content-Type: application/json` で送信してください。

---

## 共通レスポンス形式

### 成功時

```json
{
  "message": "OK",
  "results": [
    {
      "operation": "telop_change",
      "success": true
    },
    {
      "operation": "scene_change",
      "success": true
    }
  ]
}
```

- `message`: 処理結果の概要
- `results`: 各内部処理（テロップ変更、シーン変更など）の成否リスト

### エラー時

```json
{
  "message": "error description",
  "code": 3
}
```

- `message`: エラーの詳細メッセージ
- `code`: アプリケーションエラーコード（詳細は後述）

エラーレスポンスの HTTP ステータスコードはエラーの種類に応じて変わります（主に `400 Bad Request` または `500 Internal Server Error`）。

---

## エンドポイント一覧

### GET /health

ヘルスチェックエンドポイントです。

#### レスポンス

```json
{
  "message": "Hello World"
}
```

---

### GET /performances

`events.json` に定義された全パフォーマンス情報を取得します。

#### レスポンス

```json
[
  {
    "id": "d2a3e5c1-4b6f-47f7-bd8e-0e1cbe7f2c11",
    "title": "title",
    "performer": "performer",
    "description": "説明",
    "starts_at": "2025-11-01T18:00:00+09:00",
    "ends_at": "2025-11-01T19:30:00+09:00",
    "musics": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "title": "title",
        "artist": "artist",
        "should_be_muted": false,
        "intro": "intro..."
      }
    ]
  }
]
```

- `starts_at` / `ends_at`: ISO 8601 (RFC3339) 形式
- `should_be_muted`: 楽曲再生時にミュートが必要かどうか

---

### POST /force_mute

強制ミュートを ON/OFF します。

強制ミュートが ON の間、自動ミュート解除が無効化され、手動で OFF にするまでミュート状態が維持されます。

#### リクエストボディ

```json
{
  "is_muted": true
}
```

#### レスポンス例（強制ミュート ON）

```json
{
  "message": "OK",
  "results": [
    { "operation": "force_mute_on", "success": true },
    { "operation": "mute_change", "success": true }
  ]
}
```

#### レスポンス例（強制ミュート OFF）

```json
{
  "message": "OK",
  "results": [
    { "operation": "force_mute_off", "success": true },
    { "operation": "mute_change", "success": true }
  ]
}
```

#### 制約

- CM シーン中は強制ミュートを ON にできません（`400` エラー）。

---

### POST /performance/start

新しいパフォーマンスを開始します。

このエンドポイントで以下が実行されます：

1. テロップを「パフォーマンス開始」情報に更新
2. WebSocket で Viewer に通知
3. OBS シーンを **Normal** に切り替え

#### リクエストボディ

```json
{
  "title": "Evening Jazz Session",
  "performer": "Blue Note Quartet"
}
```

#### レスポンス例

```json
{
  "message": "OK",
  "results": [
    { "operation": "telop_change", "success": true },
    { "operation": "scene_change", "success": true }
  ]
}
```

---

### POST /performance/music

パフォーマンス中の楽曲を切り替えます。

このエンドポイントで以下が実行されます：

1. テロップを楽曲情報に更新
2. WebSocket で Viewer に通知
3. `should_be_muted` に応じて OBS シーンを **Muted** / **Normal** に自動切替

#### リクエストボディ

```json
{
  "title": "Autumn Leaves",
  "artist": "Joseph Kosma",
  "should_be_muted": false
}
```

#### レスポンス例

```json
{
  "message": "OK",
  "results": [
    { "operation": "telop_change", "success": true },
    { "operation": "mute_change", "success": true }
  ]
}
```

#### 注意

- 強制ミュート中 (`force_mute`) はミュート状態を変更できません。

---

### POST /conversion/start

転換（次のパフォーマンスへの準備時間）を開始します。

このエンドポイントで以下が実行されます：

1. テロップを「次のパフォーマンス情報」に更新
2. WebSocket で Viewer に通知
3. 強制ミュート状態に応じて OBS シーンを **Normal** または **Muted** に切り替え

#### リクエストボディ

```json
{
  "next_performances": [
    {
      "title": "Rock Night",
      "performer": "The Wild Guitars",
      "description": "熱狂的なロックナイト！",
      "starts_at": "2025-11-01T20:00:00+09:00"
    }
  ]
}
```

- `starts_at`: ISO 8601 (RFC3339) 形式

#### レスポンス例

```json
{
  "message": "OK",
  "results": [
    { "operation": "telop_change", "success": true },
    { "operation": "Normal_Scene_change", "success": true }
  ]
}
```

---

### POST /conversion/cm-mode

転換中の CM モードを ON/OFF します。

このエンドポイントは **転換パート中のみ有効** です。

- `is_cm_mode: true` → OBS シーンを **CM** に切り替え
- `is_cm_mode: false` → OBS シーンを **Normal** または **Muted**（強制ミュート中）に戻す

#### リクエストボディ

```json
{
  "is_cm_mode": true
}
```

#### レスポンス例

```json
{
  "message": "OK",
  "results": [
    { "operation": "CM_Scene_change", "success": true },
    { "operation": "telop_change", "success": true }
  ]
}
```

#### エラー

- 転換パート以外で呼び出した場合は `400 Bad Request` を返します。

---

### POST /display-copyright

著作権表示の ON/OFF を Viewer に通知します。

OBS シーンには影響しません。WebSocket で通知されるのみです。

#### リクエストボディ

```json
{
  "is_displayed_copyright": true
}
```

#### レスポンス例

```json
{
  "message": "OK",
  "results": [
    { "operation": "telop_change", "success": true }
  ]
}
```

---

### GET /ws

Viewer 用 WebSocket 接続エンドポイントです。

接続後、サーバーからリアルタイムにテロップ情報を受信できます。詳細は `WEBSOCKET.md` を参照してください。

---

## エラーコード一覧

| コード | 定数名 | 説明 | HTTP ステータス |
|--------|--------|------|----------------|
| `0` | `NoConnectionToOBS` | OBS への接続に失敗 | `500` |
| `1` | `NoConnectionToViewer` | Viewer への接続に失敗 | `500` |
| `2` | `InvalidPerformancesJson` | `events.json` の読み込み / パース失敗 | `500` |
| `3` | `InvalidFormat` | リクエストの形式が不正 | `400` |
| `4` | `CannotConversion` | 転換処理に失敗 | `400` |
| `5` | `CannotForceMute` | 強制ミュート処理に失敗 | `400` |
| `6` | `CannotChangeState` | 状態変更が許可されていない | `500` |

---

## Swagger UI

本 API の全スキーマ定義は、サーバー起動後に以下で確認できます：

```
http://localhost:8080/swagger/index.html
```
