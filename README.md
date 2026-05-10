# DDD_A-to-Z

DDDハッカソン用プロジェクト「Lang War」のリポジトリです。

## Quick Start

推奨環境は Nix + direnv です。

```bash
direnv allow
pnpm install
pnpm dev
```

direnv を使わない場合は、以下でも同じ開発環境に入れます。

```bash
nix develop
pnpm install
pnpm dev
```

起動後の URL は以下です。

- Web: http://localhost:5173
- API health check: http://localhost:8080/healthz

5173 が使用中の場合、Vite が表示する代替ポートを使ってください。

## Commands

```bash
pnpm dev
pnpm lint
pnpm test
pnpm build
pnpm format
```

ルートのコマンドは各言語・各 workspace の検証を呼び出すための入口です。TypeScript の整形は `apps/web` の Prettier、Go の整形は `gofmt` が担当します。
monorepo のタスク実行は Turbo で管理します。

## Requirements

Nix を使う場合、Node.js / pnpm / Go は `flake.nix` で提供されます。
Nix を使わない場合は、以下をローカルに用意してください。

- Node.js 22
- pnpm 10
- Go 1.26

## Project Layout

```text
apps/
  web/          React + Vite の最小フロントエンド
services/
  api/          Go の最小 HTTP API
docs/           企画・技術・アーキテクチャ資料
```

現時点ではプロジェクト本体の機能は実装していません。新規メンバーが同じ環境で開発を始められるようにするための基盤だけを置いています。

## Docs

- [ゲーム企画書](docs/ProductOverview.md)
- [技術設計書](docs/Tech.md)
- [アーキテクチャ設計書](docs/Architecture.md)
- [Turbo 運用ガイド](docs/Turbo.md)
