# Lang War Dlang Load Tester

D言語（Dlang）で実装した、Lang War API 向けの defensive load testing CLI です。

このツールは攻撃用途ではなく、許可された自サービス環境で API の過負荷時挙動、timeout、status 分布、latency を確認するための開発者向け補助ツールです。

## Safety

- デフォルト target は `http://localhost:8080/healthz` です。
- localhost 以外の target は `--allow-external` を付けない限り実行できません。
- `--concurrency` は `200` までに制限しています。
- `--duration` は `1h` まで、`--timeout` は `5m` までに制限しています。
- 許可されたローカル環境、staging 環境、自分たちが管理する環境にだけ実行してください。
- CI では高負荷テストを実行せず、build / unit test のみを行います。

## Usage

```bash
pnpm --filter @lang-war/loadtest-d build
pnpm --filter @lang-war/loadtest-d test
```

`dub.json` も置いているため、DUB で直接扱うこともできます。

```bash
dub build --compiler=ldc2
dub test --compiler=ldc2
```

低負荷のローカル実行例:

```bash
pnpm --filter @lang-war/loadtest-d dev -- \
  --target http://localhost:8080/healthz \
  --duration 10s \
  --concurrency 5 \
  --timeout 2s \
  --output json
```

外部の staging URL に実行する場合は明示的に許可します。

```bash
pnpm --filter @lang-war/loadtest-d dev -- \
  --target https://example-staging.example.com/healthz \
  --duration 10s \
  --concurrency 5 \
  --timeout 2s \
  --output json \
  --allow-external
```

## Output

`successRequests` は HTTP 200-399（2xx-3xx）の response を成功として数えます。
Redirect は追跡しないため、3xx response は 3xx のまま `statusCounts` に記録され、成功扱いになります。

```json
{
  "target": "http://localhost:8080/healthz",
  "durationMs": 10000,
  "concurrency": 5,
  "totalRequests": 1200,
  "successRequests": 1200,
  "failedRequests": 0,
  "statusCounts": {
    "200": 1200
  },
  "errorCounts": {},
  "latencyMs": {
    "avg": 12.4,
    "p95": 31,
    "p99": 49
  }
}
```

## Options

| Option | Default | Description |
|---|---:|---|
| `--target` | `http://localhost:8080/healthz` | GET request の送信先 |
| `--duration` | `5s` | 実行時間。`ms`, `s`, `m`, `h` を指定可能 |
| `--concurrency` | `2` | worker thread 数 |
| `--timeout` | `2s` | request ごとの timeout |
| `--output` | `json` | `json` または `text` |
| `--allow-external` | `false` | localhost 以外への実行を明示的に許可 |

成功判定は HTTP 200-399 です。Redirect は追跡しないため、3xx response もそのまま成功として集計されます。
