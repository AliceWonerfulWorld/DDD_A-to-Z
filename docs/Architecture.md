# アーキテクチャ設計書 — Lang War

## 方針

Lang War は、MVP では単一の Go API と React フロントエンドで始める。最初に作るのは「GitHub活動が言語別ポイントに変換され、ランキングが動く」体験であり、マイクロサービス構成そのものではない。

将来的には A〜Z の各技術・言語を独立した Cloud Run サービスとして追加し、ゲーム内の言語勢力と実際のサービス構成が対応する形を目指す。ただし、その分割はコア体験が成立してから行う。

## MVP 全体構成

```text
[React Frontend]
  Vercel / Cloudflare Pages / Firebase Hosting
        |
        v
[Go API]
        |
        +--> [Cloud SQL: PostgreSQL]
        |
        +--> [GitHub API]
```

## コンポーネント

| コンポーネント | 役割 | スケール方針 |
|---|---|---|
| `frontend` | React/Vite の画面 | 静的ホスティング |
| `api-go` | 画面 API、認証、GitHub 連携、ポイント変換、ランキング表示 | 最初は単一プロセス |
| `postgres` | ユーザー、言語、シーズン、GitHub活動、ポイント履歴 | マネージド DB |

## リクエストフロー

### ユーザー画面

ユーザー向け画面では、ランキングと自分の貢献ログを PostgreSQL の保存済みデータから返す。

```text
React
  -> api-go
    -> PostgreSQL からランキング・貢献ログを取得
    -> 画面用レスポンスを返す
```

### GitHub活動取り込み

MVPでは常時Webhook受信を前提にしない。ユーザーがGitHub連携を行ったあと、画面操作または定期的な軽い同期で直近活動を取得する。

```text
React
  -> api-go
    -> GitHub REST API から直近活動を取得
    -> 言語判定
    -> ポイント計算
    -> PostgreSQL に保存
```

## フロントエンド API クライアント構成

フロントエンドは2種類の API クライアントを使い分ける。どちらも `apps/web/src/lib/api/` に集約している。

```text
apps/web/src/lib/api/
├── client.ts    — HTTP JSON (REST) クライアント。apiFetch() を提供
└── connect.ts   — Connect RPC クライアント。connectClient() を提供
```

### 使い分け

| 通信種別 | 関数 | 用途 |
|---|---|---|
| HTTP JSON (REST) | `apiFetch<T>(path, init?)` | 認証・OAuth・未変換のエンドポイント |
| Connect RPC | `connectClient(Service)` | proto 定義済みのエンドポイント |

### Connect RPC クライアントの使い方

`connectClient()` は transport をシングルトンで保持し、任意のサービスのクライアントを返す。transport の設定（interceptor 追加、baseUrl 変更など）は `connect.ts` の1箇所で管理する。

```typescript
// features/home/api.connect.ts
import { connectClient } from "../../lib/api/connect";
import { HomeService } from "@lang-war/proto-ts/langwar/home/v1/home_service_pb";

const client = connectClient(HomeService);

export async function fetchHomeViaConnect() {
  return client.getHome({});
}
```

新しい proto サービスを追加した場合も `connectClient(NewService)` の1行で済む。

### 共存方針

REST と Connect RPC は同一ポートで共存している。新しいエンドポイントは最初から proto で定義し、既存の REST エンドポイントは段階的に移行する。認証・OAuth フローは REST のまま維持する。

## スケール方針

MVPではスケールよりも検証速度を優先する。ランキング集計はPostgreSQLで行い、更新頻度やレスポンス時間が問題になってからキャッシュや非同期処理を追加する。

| 対象 | 方針 |
|---|---|
| `api-go` | まずは1サービス。必要になったらBFF/workerに分ける |
| `postgres` | 最初のランキング集計を担当する |
| キャッシュ | PostgreSQL集計が遅くなってから導入 |
| 非同期処理 | GitHub同期やシーズン締めが重くなってから導入 |

## MVP 構成

MVP では、言語別サービスを最初から 26 個作らない。まずは単一の Go API の中にドメインロジックを置き、GitHub活動を言語別ポイントに変換してランキングとして表示する。

MVPのコア体験は「GitHubで書いたコードが、自分の推し言語の勢力ポイントになり、ランキングが動くのを見る」こと。テリトリーマップ、キャラクター能力値、リアルタイムWebhook処理、多言語サービス分割は後回しにする。

```text
frontend
api-go
postgres
```

初期実装では、画面からの操作でGitHub REST APIを呼び、直近活動を取り込む。ランキング集計はPostgreSQL上で行い、負荷や更新頻度が問題になってから非同期worker、Valkey、queueを追加する。

## 拡張後の構成

主要言語から順に Cloud Run サービスとして追加する。

```text
web-bff-go
worker-go
lang-service-go
lang-service-python
lang-service-rust
lang-service-java
...
```

最終的には、A〜Z の 26 言語・技術サービスに BFF と worker を加えて、Cloud Run サービスが 30 個弱立つ構成を想定する。

## ファイル構成

MVP は monorepo とし、フロントエンド、単一の Go API、DB、ローカル開発用 infra を同じリポジトリで管理する。

```text
DDD_A-to-Z/
├── README.md
├── docs/
│   ├── ProductOverview.md
│   ├── Tech.md
│   ├── Architecture.md
│   └── Deployment.md
│
├── apps/
│   └── web/
│       ├── package.json
│       ├── vite.config.ts
│       ├── src/
│       │   ├── main.tsx
│       │   ├── app/
│       │   ├── pages/
│       │   ├── features/
│       │   ├── components/
│       │   └── lib/
│       └── Dockerfile
│
├── proto/                        # Protobuf スキーマ定義
│   ├── buf.yaml
│   ├── buf.gen.go.yaml
│   ├── buf.gen.web.yaml
│   └── langwar/
│       └── home/
│           └── v1/
│
├── gen/                          # buf generate で生成されるコード (コミット対象)
│   ├── go/                       # Go 生成コード (go.work 独立モジュール)
│   └── ts/                       # TypeScript 生成コード (@lang-war/proto-ts パッケージ)
│
├── services/
│   └── api/
│       ├── cmd/
│       │   └── server/
│       │       └── main.go
│       ├── internal/
│       │   ├── domain/
│       │   ├── application/
│       │   ├── infrastructure/
│       │   └── interfaces/
│       ├── go.mod
│       └── Dockerfile
│
├── packages/
│   └── shared-types/
│
├── db/
│   ├── migrations/
│   ├── schema.sql
│   └── seed/
│
├── infra/
│   ├── docker-compose.yml
│   ├── cloudrun/
│   └── terraform/
│
├── scripts/
│   ├── atlas.sh
│   └── buf.sh
│
└── .github/
    └── workflows/
        ├── ci-api.yml
        ├── ci-proto.yml
        └── deploy.yml
```

## 将来のサービス追加イメージ

A〜Z の各技術を増やす場合は、`services/` 配下にサービスを追加する。

```text
services/
├── web-bff/
├── worker/
├── lang-astro/
├── lang-bun/
├── lang-cpp/
├── lang-deno/
├── lang-elixir/
├── lang-fsharp/
├── lang-go/
├── lang-haskell/
├── lang-java/
├── lang-kotlin/
├── lang-python/
├── lang-rust/
├── lang-typescript/
└── lang-zig/
```

## 設計上の注意点

- MVP では Go 内部でモジュール分割し、外部サービス化が必要になった処理から切り出す。
- ユーザー導線上で 26 サービスを同期呼び出ししない。
- `lang-service-*` を追加する場合も、DB に直接書き込ませない。
- `proto/` は全サービス共通の契約として導入済み。新しい言語サービスを追加する際は同じ `.proto` を実装するだけで接続できる。
- Cloud Run のゼロスケールは、ユーザー導線から外れた非同期サービスを中心に使う。

## DDD 方針

Lang War では、技術を増やすこと自体を目的にせず、ゲームの中核ルールをドメインモデルとして明確に保つ。MVP では単一の Go サービス内に実装してよいが、コード上はドメイン、ユースケース、インフラを分ける。

### 境界づけられたコンテキスト

初期段階では、以下のコンテキストを意識する。

| コンテキスト | 役割 |
|---|---|
| Game / Faction | 言語勢力、ポイント、ランキング、シーズンの中核ルール |
| GitHub Integration | GitHub OAuth、GitHub REST API、活動データの正規化 |
| User / Account | ユーザー、GitHub アカウント連携、認証状態 |
| Presentation / API | 画面向け API、表示用集約、レスポンス整形 |
| Language Runtime | MVP後の言語別サービスによる計算、評価、キャラクター能力値算出 |

最初からこれらを物理サービスとして分割しない。まずは同一 Go プロセス内の package 境界として分け、負荷・変更頻度・チーム分担が見えてから外部サービス化する。

### レイヤー

Go サービス内では、以下の依存方向を守る。

```text
cmd
  -> http handlers
    -> application
      -> domain
    -> infra
```

`domain` は GitHub API や PostgreSQL などの外部 I/O を直接知らない。外部 I/O は `application` から interface 経由で呼び、実装は `infra` 側に置く。

### 集約候補

MVP で中心になる集約は以下。

| 集約 | 説明 |
|---|---|
| `Season` | シーズン期間、状態、開始・終了ルール |
| `Faction` | 言語勢力。現在ポイント、順位、キャラクター状態を持つ |
| `GitHubEvent` | GitHub から受け取った活動イベント。重複排除と処理状態を持つ |
| `Contribution` | ユーザー活動をゲーム上のポイントに変換した結果 |
| `UserAccount` | ユーザーと GitHub アカウント連携 |

ポイント計算やシーズン締めなど、ゲームルールとして重要な処理は `domain` に寄せる。DB の都合や GitHub の payload 形式をそのままドメインモデルに持ち込まない。

### ユースケース候補

`application` 層には、画面や同期処理から呼ばれるユースケースを置く。

| ユースケース | 呼び出し元 | 概要 |
|---|---|---|
| `GetRanking` | API | PostgreSQL からランキングを取得する |
| `GetMyContributions` | API | ユーザーのGitHub活動とポイント変換結果を取得する |
| `SyncGitHubActivity` | API | GitHub REST APIから直近活動を取得し、ポイントへ変換する |
| `GetLanguageDetail` | API | 言語ごとのポイント、順位、簡単なキャラクター情報を取得する |
| `CloseSeason` | admin / batch | シーズンを締め、結果を確定する。MVP後に自動化する |

### Repository と外部サービス

ドメインルールから永続化や外部 API を直接呼ばない。`application` 層が必要な port interface を定義し、GitHub client や PostgreSQL repository は `infrastructure` 側で実装する。

```text
interfaces/http
  -> application use case
    -> domain

application ports
  -> FactionRepository
  -> SeasonRepository
  -> GitHubEventRepository
  -> ContributionRepository
  -> GitHubActivityClient

infrastructure
  -> postgres repository
  -> github client
```

### パッケージ構成イメージ

MVPでは `services/api` の中で package 境界を分ける。重複がつらくなったら `services/internal` や `packages/domain-go` のような共有 package を検討する。ただし、早すぎる共有化は避ける。

```text
services/api/internal/
├── domain/
│   ├── user/
│   ├── repositoryanalysis/
│   ├── contributionpoint/
│   └── guild/
├── application/
│   ├── auth/
│   ├── user/
│   ├── github/
│   ├── repositoryanalysis/
│   ├── contributionpoint/
│   ├── guild/
│   └── mypage/
├── infrastructure/
│   ├── config/
│   ├── database/
│   ├── github/
│   ├── postgres/
│   ├── memory/
│   ├── security/
│   └── logger/
├── interfaces/
│   └── http/
```

### マイクロサービス化の基準

DDD の境界を、そのまま最初から Cloud Run サービス境界にしない。外部サービス化は、次の条件が出てから行う。

- 起動・ランタイム・依存関係が明確に違う。
- 負荷特性が API 本体と違う。
- 変更頻度や担当者が分かれる。
- 失敗しても本体のユーザー導線に影響させたくない。
- 言語別サービスとして独立させること自体がゲーム体験になる。

この方針により、MVP では DDD のコード境界を保ちつつ、運用上のマイクロサービス分割は後回しにする。

## 開発環境の方針

開発中のツール管理は、`pnpm`、Nix flakes、Docker Compose を役割ごとに分ける。

```text
pnpm           JS/TS workspace、root scripts、frontend
Go tooling     services/api のビルド・テスト
Nix flakes     Node、pnpm、Go、Docker 系 CLI などの開発ツール固定
Docker Compose ローカルの PostgreSQL
Atlas          DB schema / migration 管理
```

`pnpm` は全体のコマンド入口として使う。ただし Go サービス自体を無理に pnpm workspace package にしない。root の `package.json` から `go run` や `go test` を呼び出す形にする。

Nix flakes は、まずは開発ツールのバージョン固定に用途を絞る。production build や各サービスのコンテナ生成まで Nix に寄せるのは、必要になってから検討する。

## 開発用ファイル構成

開発ツールを入れた root の構成は以下を想定する。

```text
DDD_A-to-Z/
├── package.json
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
├── flake.nix
├── flake.lock
├── .env.example
├── apps/
│   └── web/
├── services/
│   └── api/
├── db/
└── infra/
    └── docker-compose.yml
```

`pnpm-workspace.yaml` は JS/TS の package を対象にする。

```yaml
packages:
  - "apps/*"
  - "gen/*"
  - "packages/*"
  - "services/*"
```

`gen/ts` は `@lang-war/proto-ts` として workspace パッケージ化されており、`apps/web` から `workspace:*` 参照している。将来言語が増えた場合も `gen/<lang>/` に同様のパッケージを追加するだけで対応できる。

## root scripts

root の `package.json` は、開発中によく使う操作の入口にする。

```json
{
  "scripts": {
    "dev:web": "pnpm --filter web dev",
    "dev:api": "go run ./services/api/cmd/server",
    "dev:infra": "docker compose -f infra/docker-compose.yml up",
    "test": "pnpm -r test && go test ./services/...",
    "lint": "pnpm -r lint"
  }
}
```

MVP 初期は `dev:all` を無理に作らなくてもよい。フロントエンド、API、infra を個別ターミナルで起動する方が挙動を追いやすい。必要になった段階で、複数プロセスをまとめて起動する `dev` script を追加する。

## flake.nix

`flake.nix` は最初は dev shell だけを定義する。

```nix
{
  description = "Lang War development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      systems = [ "aarch64-darwin" "x86_64-darwin" "x86_64-linux" "aarch64-linux" ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs systems (system:
          f nixpkgs.legacyPackages.${system}
        );
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [
            pkgs.nodejs
            pkgs.corepack
            pkgs.go
            pkgs.docker
            pkgs.docker-compose
            pkgs.git
            pkgs.ripgrep
          ];

          shellHook = ''
            corepack enable
          '';
        };
      });
    };
}
```

Nix flakes は実務ではよく使われるが、Nix 公式上は experimental 機能として扱われている。そのため、初期段階では「開発環境の固定」に限定して使う。

## ローカル開発フロー

基本の流れは以下。

```bash
nix develop
pnpm install
pnpm dev:infra
pnpm dev:web
pnpm dev:api
```

ローカルの PostgreSQL は Docker Compose で起動する。DB schema は Go API の起動処理ではなく Atlas で管理し、`db/schema.sql` を source of truth、`db/migrations/` を versioned migration 置き場にする。アプリケーションはホスト側で起動し、ホットリロードやデバッガを使いやすくする。

## まだ入れないもの

初期段階では、開発基盤を重くしすぎない。

| ツール | 判断 |
|---|---|
| Turbo | package が増えて build/test キャッシュが必要になってから導入 |
| Nx | 今の規模では重い |
| devenv | 強力だが、Nix flakes + pnpm + Docker Compose で足りる間は入れない |
| production build の Nix 化 | Cloud Run 用 Dockerfile が安定してから検討 |
| Valkey | PostgreSQL集計で足りなくなってから導入 |
| Pub/Sub / Cloud Tasks | 非同期処理が必要になってから導入 |
| gRPC / 他プロトコル移行 | Connect RPC (HTTP/1.1 + JSON) で足りる間は不要 |
| 言語別サービス | Go単体でコア体験を検証してから分離 |

当面の採用方針は、`pnpm workspace`、`flake.nix`、Docker Compose、Go 標準 tooling とする。
