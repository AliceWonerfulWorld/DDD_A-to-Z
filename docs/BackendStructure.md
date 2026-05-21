# バックエンドファイル構成

MVP の Go API は DDD と Clean Architecture を基本方針にする。内側ほどビジネスルールに近く、外側ほど技術詳細に近い。

依存方向は常に以下にそろえる。

```text
interfaces -> application -> domain
infrastructure -> application/domain
cmd/server -> 全体の組み立て
```

`domain` は GitHub API、PostgreSQL、HTTP、env を知らない。`application` は use case と port interface を持ち、`infrastructure` がその port を実装する。

## ディレクトリ

```text
services/api/
├── cmd/
│   └── server/                       # DI、プロセス起動、ListenAndServe
└── internal/
    ├── domain/                       # Entity / Value Object / Domain Service
    │   ├── user/                     # User、GitHubAccount
    │   ├── repositoryanalysis/       # リポジトリ解析のドメイン概念
    │   ├── contributionpoint/                       # ContributionPoint のルール
    │   └── guild/                    # Guild、ランキング、所属状態
    ├── application/                  # Use Case / Port
    │   ├── auth/                     # GitHubログイン use case
    │   ├── user/
    │   ├── github/
    │   ├── repositoryanalysis/
    │   ├── contributionpoint/
    │   ├── guild/
    │   └── mypage/
    ├── infrastructure/               # 外部I/Oの実装
    │   ├── config/                   # env、設定値
    │   ├── database/                 # DB接続、transaction helper
    │   ├── github/                   # GitHub OAuth / REST client
    │   ├── postgres/                 # PostgreSQL repository 実装
    │   ├── memory/                   # 開発・テスト用 in-memory 実装
    │   ├── security/                 # token生成など技術詳細
    │   └── logger/                   # logger初期化
    ├── interfaces/
    │   ├── http/                     # controller、router、request/response変換
    │   └── connect/                  # Connect RPC ハンドラ (Protobuf)
```

## フォルダごとの説明

### `cmd/server`

API サーバープロセスの入口。ここでは logger、config、repository、use case、controller を組み立てて `ListenAndServe` する。

置いてよいもの:

- `main.go`
- DI の組み立て
- プロセス起動時の logger / config 読み込み

置かないもの:

- ユーザー登録やログインなどの業務ロジック
- SQL
- GitHub API 呼び出し
- HTTP request の細かい validation

### `internal/domain`

アプリの中核ルールを置く場所。DDD でいう Entity、Value Object、Domain Service を置く。ここは一番内側なので、外部技術を知らない。

置いてよいもの:

- `User`, `GitHubAccount`, `Guild` などの Entity
- ID、ポイント、期間などの Value Object
- 「GitHub ID は必須」「username は空にできない」のような不変条件
- DB や API に依存しない計算ルール

置かないもの:

- `net/http`
- SQL
- env 読み込み
- GitHub API client
- JSON request / response 用の都合

例:

- `domain/user`: アプリ内 User と GitHubAccount の関係
- `domain/contributionpoint`: ContributionPoint の計算ルール
- `domain/guild`: 言語ギルド、ランキング、所属状態のルール
- `domain/repositoryanalysis`: リポジトリ解析結果のドメイン概念

### `internal/application`

ユースケースを置く場所。画面やバッチから呼ばれる「アプリとして何をするか」を表す。

置いてよいもの:

- `BeginGitHubLogin`
- `CompleteGitHubLogin`
- `GetMyPage`
- `SyncGitHubActivity`
- Repository / Client の interface
- transaction 境界の制御

置かないもの:

- SQL の具体実装
- `http.ResponseWriter`
- GitHub REST API の URL や payload 詳細
- PostgreSQL driver 固有の処理

ポイント:

- `application` は `domain` を使って処理する。
- DB や外部 API が必要な場合は、ここに interface を定義する。
- 実装は `infrastructure` に置く。

例:

```text
application/auth/usecase.go      # ログイン use case
application/auth/ports.go        # GitHubOAuthClient, UserRepository などの interface
```

### `internal/infrastructure`

外部 I/O や技術詳細を置く場所。`application` で定義した interface を実装する。

置いてよいもの:

- GitHub OAuth / REST API client
- PostgreSQL repository
- DB 接続、transaction helper
- env 読み込み
- token 生成
- logger 初期化
- 開発・テスト用 in-memory repository

置かないもの:

- HTTP controller
- use case の判断ロジック
- domain の不変条件そのもの

例:

- `infrastructure/github`: GitHub OAuth / REST API の実装
- `infrastructure/postgres`: PostgreSQL repository の実装
- `infrastructure/memory`: 開発・テスト用の in-memory 実装
- `infrastructure/config`: env から設定を読む
- `infrastructure/database`: DB 接続や transaction helper
- `infrastructure/security`: token 生成など
- `infrastructure/logger`: logger 初期化

### `internal/interfaces`

外部からアプリに入ってくる入口を置く場所。HTTP API (`interfaces/http`) と Connect RPC (`interfaces/connect`) の2つがある。

置いてよいもの:

- router
- controller / handler
- request DTO
- response DTO
- HTTP status code の選択
- cookie や header の設定
- Connect RPC ハンドラ (proto 生成インターフェースを実装)

置かないもの:

- domain ルール
- SQL
- GitHub API 呼び出し
- DB transaction の詳細

ポイント:

- controller は request を読み、use case を呼び、response を返す。
- controller が直接 repository や GitHub client を呼ばない。
- request / response の形は domain model と分けてよい。
- `interfaces/connect` は `interfaces/http` と同じ application 層の依存 (port interface) を再利用する。HTTP/JSON と Connect RPC は同一 ServeMux で共存する。

## 実装例

GitHub ログインを最初のサンプル実装として置いている。

```text
interfaces/http/auth_controller.go
  -> 署名付き cookie で OAuth state を検証
  -> application/auth/usecase.go
    -> domain/user/user.go
    -> application/auth/ports.go
      <- infrastructure/github/oauth.go
      <- infrastructure/postgres/auth_repositories.go
      <- infrastructure/security/token.go
```

この形を後続 Issue の基準にする。サーバーレス環境でも login callback が別インスタンスに届いて壊れないように、OAuth state は署名付き cookie、User と session は PostgreSQL に保存する。DB schema / migration は Go の起動処理ではなく root の Atlas 設定で管理する。

## 配置ルール

- `cmd/server` は組み立てだけを持つ。業務ロジックは置かない。
- Entity / Value Object / Domain Service は `domain` に置く。
- Use Case と Repository / Client の interface は `application` に置く。
- GitHub API、PostgreSQL、token生成などの技術詳細は `infrastructure` に置く。
- HTTP controller は `interfaces/http` に置き、use case を呼ぶだけにする。
- config、DB接続、logger なども外側の技術詳細として `infrastructure` に置く。
