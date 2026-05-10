# Turbo 運用ガイド

## 方針

このリポジトリでは、Turbo を pnpm monorepo の司令塔として使う。

ただし、Turbo に各言語のビルドや依存関係を深く理解させるのではなく、各 app / service / package が持つネイティブツールを呼ぶ薄い入口として扱う。

```text
root package.json
  -> turbo run build / test / lint / format / dev
    -> apps/web/package.json
      -> vite / tsc / prettier
    -> services/api/package.json
      -> go / gofmt
    -> future service package.json
      -> cargo / python / zig / bun / ...
```

## 基本ルール

- root は `pnpm` と `turbo` による全体実行だけを担当する。
- 各 app / service / package には最小の `package.json` を置く。
- `package.json` は Node の依存管理だけでなく、monorepo タスクのアダプターとして使う。
- 各言語の実処理は、その言語の標準ツールに任せる。
- 複雑な言語別ロジックを root scripts に集めない。

## 標準スクリプト

新しい app / service / package は、可能な限り以下の script を持つ。

```json
{
  "scripts": {
    "dev": "...",
    "build": "...",
    "test": "...",
    "lint": "...",
    "format": "..."
  }
}
```

全ての package が必ず `dev` を持つ必要はない。ライブラリやバッチのように常駐起動がない場合は、`build` / `test` / `lint` / `format` を優先する。

## 追加手順

新しい package を追加するときは、以下の流れにする。

1. `apps/*`, `services/*`, `packages/*` のどこに置くか決める。
2. そのディレクトリに `package.json` を置く。
3. `name` は `@lang-war/...` 形式にする。
4. 標準スクリプトを、その言語のネイティブツールを呼ぶ形で定義する。
5. `pnpm install` で lockfile を更新する。
6. `pnpm lint`, `pnpm test`, `pnpm build` が root から通ることを確認する。

`pnpm-workspace.yaml` は現在 `apps/*`, `services/*`, `packages/*` を対象にしているため、この配下に置けば Turbo の対象に入る。

## 例: Go service

```json
{
  "name": "@lang-war/api",
  "private": true,
  "scripts": {
    "dev": "GOCACHE=$(pwd)/.cache/go-build go run ./cmd/server",
    "build": "mkdir -p .cache/bin && GOCACHE=$(pwd)/.cache/go-build go build -o .cache/bin/lang-war-api ./cmd/server",
    "test": "GOCACHE=$(pwd)/.cache/go-build go test ./...",
    "lint": "test -z \"$(gofmt -l cmd internal)\"",
    "format": "gofmt -w cmd internal"
  }
}
```

## 例: Rust service

```json
{
  "name": "@lang-war/lang-rust",
  "private": true,
  "scripts": {
    "build": "cargo build",
    "test": "cargo test",
    "lint": "cargo clippy -- -D warnings",
    "format": "cargo fmt"
  }
}
```

## 例: TypeScript package

```json
{
  "name": "@lang-war/shared-types",
  "private": true,
  "type": "module",
  "scripts": {
    "build": "tsc -b",
    "test": "tsc -b",
    "lint": "pnpm format:check && tsc -b",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
}
```

## キャッシュの育て方

Turbo は Go や Rust の依存関係を深く理解するわけではない。

最初は `turbo.json` の共通設定だけでよいが、package が増えて build 時間が気になり始めたら、タスクごとに `inputs` / `outputs` を明示する。

例。

```json
{
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".cache/bin/**", "target/**"]
    }
  }
}
```

キャッシュ設定は、実際の言語と成果物が増えてから育てる。最初から Bazel や Pants のような重量級ルールを目指さず、ハッカソンで参加しやすい薄い規約を優先する。
