#!/bin/sh
set -e

# env bundle secret をファイルマウントから環境変数へ展開
if [ -f /run/secrets/runtime-env ]; then
  export $(grep -v '^#' /run/secrets/runtime-env | xargs)
fi

# APP_KEY が未設定の場合はエラー
if [ -z "$APP_KEY" ]; then
  echo "ERROR: APP_KEY is not set" >&2
  exit 1
fi

exec "$@"
