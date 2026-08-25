# WebSocket 仕様書

## 接続

Viewer は以下のエンドポイントに WebSocket 接続します。

```
ws://{host}:8080/ws
```

接続元（Origin）は `localhost`, `127.0.0.1`, `::1` のみ許可されています。

---

## メッセージ形式

サーバーから送信されるメッセージはすべて **テキストフレーム（JSON）** です。

```json
{
  "type": "/performance/start",
  "data": { ... }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `type` | `string` | メッセージの種別（エンドポイントのパスと対応） |
| `data` | `object` | 種別に応じたペイロード |

---

## メッセージ種別

### 1. `/performance/start`

パフォーマンスが開始されたときに送信されます。

#### ペイロード (`PerformanceStartData`)

```json
{
  "type": "/performance/start",
  "data": {
    "title": "title",
    "performer": "performer"
  }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `title` | `string` | パフォーマンスタイトル |
| `performer` | `string` | 出演者名 |

---

### 2. `/performance/music`

パフォーマンス中の楽曲が切り替わったときに送信されます。

#### ペイロード (`PerformanceMusicData`)

```json
{
  "type": "/performance/music",
  "data": {
    "title": "title",
    "artist": "artist",
    "should_be_muted": false
  }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `title` | `string` | 楽曲タイトル |
| `artist` | `string` | アーティスト名 |
| `should_be_muted` | `boolean` | ミュートが必要かどうか |

---

### 3. `/conversion/start`

転換（次のパフォーマンスへの準備時間）が開始されたときに送信されます。

#### ペイロード (`ConversionStartData`)

```json
{
  "type": "/conversion/start",
  "data": {
    "next_performances": [
      {
        "title": "title",
        "performer": "performer",
        "description": "description",
        "starts_at": "2025-11-01T20:00:00+09:00"
      }
    ]
  }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `next_performances` | `array` | 次のパフォーマンス情報リスト |
| `next_performances[].title` | `string` | タイトル |
| `next_performances[].performer` | `string` | 出演者名 |
| `next_performances[].description` | `string` | 説明文 |
| `next_performances[].starts_at` | `string` | 開始時刻（ISO 8601 / RFC3339） |

---

### 4. `/conversion/cm-mode`

転換中に CM モードが ON/OFF されたときに送信されます。

#### ペイロード (`ConversionCmModeData`)

```json
{
  "type": "/conversion/cm-mode",
  "data": {
    "is_cm_mode": true
  }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `is_cm_mode` | `boolean` | CM モードが有効かどうか |

---

### 5. `/display-copyright`

著作権表示の ON/OFF が切り替えられたときに送信されます。

#### ペイロード (`DisplayCopyrightData`)

```json
{
  "type": "/display-copyright",
  "data": {
    "is_displayed_copyright": true
  }
}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `is_displayed_copyright` | `boolean` | 著作権表示が有効かどうか |

---

## クライアント実装例（JavaScript）

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to ultradonguri server');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received type:', message.type);
  console.log('Data:', message.data);

  switch (message.type) {
    case '/performance/start':
      // パフォーマンス開始時の処理
      break;
    case '/performance/music':
      // 楽曲切替時の処理
      break;
    case '/conversion/start':
      // 転換開始時の処理
      break;
    case '/conversion/cm-mode':
      // CM モード切替時の処理
      break;
    case '/display-copyright':
      // 著作権表示切替時の処理
      break;
  }
};

ws.onerror = (err) => {
  console.error('WebSocket error:', err);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```
