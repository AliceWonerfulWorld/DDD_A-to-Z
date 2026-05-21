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

if [ "${1:-}" = "breaking" ]; then
  against_ref="${3:-}"
  if [ "$against_ref" = "../.git#branch=origin/main,subdir=proto" ]; then
    if ! proto_files="$(git -C "$ROOT" ls-tree -r --name-only origin/main -- proto 2>/dev/null)"; then
      echo "Failed to inspect origin/main for proto files." >&2
      exit 1
    fi
    if ! printf '%s\n' "$proto_files" | grep -q '\.proto$'; then
      echo "No proto files on origin/main; skipping breaking check."
      exit 0
    fi
  fi
fi

exec buf "$@"
