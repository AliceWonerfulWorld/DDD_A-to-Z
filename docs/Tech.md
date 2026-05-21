# 技術設計書 — Lang War

## 技術スタック（MVP）

| レイヤー | 技術 | 選定理由 |
|---|---|---|
| フロントエンド | React + Vite | 決定済み |
| 通信 (REST) | HTTP JSON | MVPでは画面とGo APIの開発速度を優先 |
| 通信 (RPC) | Protobuf + Connect RPC | 多言語サービス化に備えた共通契約。HTTP/JSON と同一ポートで共存 |
| スキーマ管理 | buf | proto lint / build / generate / breaking change 検出 |
| バックエンド (MVP) | Go | GitHub連携と集計処理をシンプルに実装できる |
| DB | PostgreSQL | 勢力データの永続化 |
| DB schema管理 | Atlas | Go以外のサービスからも使える言語非依存の schema / migration 管理 |
| GitHub連携 | OAuth + REST API | まずはユーザー連携と直近活動の取り込みに絞る |

Valkey、Webhook常時受信は、MVP後に必要性が見えた段階で導入する。

## A〜Z 技術割り当て（たたき台）

| | 技術 | 役割 |
|---|---|---|
| A | Astro | 静的ランディングページ |
| B | Bun | ツールチェーン / スクリプト実行 |
| C | C++ | コアロジック（WASM経由） |
| D | Deno | 補助CLIツール |
| E | Elixir | リアルタイムWebSocket通知 |
| F | F# | 関数型バックエンドサービス |
| G | Go | MVPバックエンド / 将来のBFF・worker |
| H | Haskell | 勢力計算ロジックサービス |
| I | 未定 | Infra定義系（Terraform / Nix） |
| J | Java | バックエンドサービスの一つ |
| K | Kotlin | JVM系サービス |
| L | Lua | キャラ行動スクリプティング |
| M | Mojo | ML・数値計算サービス |
| N | Nim | 軽量高速サービス |
| O | OCaml | 型安全データパイプライン |
| P | Python | データ分析・グラフ生成 |
| Q | 未定 | Queue系ミドルウェア（RabbitMQ等） |
| R | Rust | 高パフォーマンスコアサービス / Wasm |
| S | Scala | ストリーム処理（Akka Streams） |
| T | TypeScript | フロントエンド本体（React） |
| U | 未定 | 難所。Unison / uv等 要検討 |
| V | Valkey | キャッシュ・リアルタイムランキング |
| W | WebAssembly | C++ / Rustのブラウザ実行 |
| X | 未定 | 最難所。xk6（負荷テスト）案あり |
| Y | Yew | RustベースのWasm UIフレームワーク |
| Z | Zig | 最低レイヤーサービス（概念上のラスボス） |

I・Q・U・Xの4つが未確定。

## フェーズ計画

**MVP** — React + Go + PostgreSQL + GitHub連携。GitHubの直近活動を言語別ポイントに変換し、シーズン内ランキングと自分の貢献ログとして表示する。

**拡張期** — MVPのGoサービスを動かしながら、Java・Python・Rustを順次バックエンドサービスとして追加していく。各サービスは同じ `.proto` を実装するだけで接続できる構成にする。`web-bff-go` が画面向けの集約を担い、`worker-go` が必要に応じて言語別サービスを呼び出す。

**最終形** — A〜Z 26技術すべてが実際にプロダクションで動いている状態。各言語のサービスがそのままゲーム内のキャラクター強度に影響する自己言及構造が完成する。

## 設計上の最重要判断

MVPのコア体験は「GitHub活動が言語勢力ランキングに反映されること」に絞った。多言語サービス化に備えて Protobuf + Connect RPC を早期導入し、HTTP/JSON と同一ポートで共存させることで既存実装を破壊せずに基盤を整えた。新しい言語サービスを追加する際は同じ `.proto` を実装するだけで接続できる。
