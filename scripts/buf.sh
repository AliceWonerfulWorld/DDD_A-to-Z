#!/usr/bin/env sh
set -eu
ROOT="$(dirname "$0")/.."
# Web (TypeScript) 生成プラグインのパスを追加
export PATH="$ROOT/apps/web/node_modules/.bin:$PATH"
cd "$ROOT/proto"

# generate は全テンプレートをまとめて実行する
# 新しい言語サービスを追加したらここに buf.gen.<lang>.yaml を追記する
if [ "${1:-}" = "generate" ]; then
  buf generate --template buf.gen.go.yaml        # Go
  buf generate --template buf.gen.web.yaml       # TypeScript (apps/web)
  exit 0
fi

exec buf "$@"
