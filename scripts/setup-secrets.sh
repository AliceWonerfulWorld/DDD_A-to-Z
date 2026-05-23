#!/usr/bin/env bash
# GCP Secret Manager にシークレットを登録するスクリプト
# terraform apply でリソースのセットアップが済んでいることが前提
# シークレットは課金を防ぐために6つに絞ってください
#
# 設定値はプロジェクトルートの .env.prod から読み込む。必須変数:
#   GCP_PROJECT_ID          — GCP プロジェクト ID
#   NEON_DATABASE_URL       — Neon DB 接続文字列
#   GITHUB_CLIENT_ID        — GitHub OAuth App Client ID
#   GITHUB_CLIENT_SECRET    — GitHub OAuth App Client Secret
#   GITHUB_REDIRECT_URL     — GitHub OAuth callback URL
#   AUTH_COOKIE_SECRET      — OAuth state cookie 署名用 secret（Chat Service の SECRET_KEY_BASE としても使うため64 bytes以上）
#   GITHUB_TOKEN_ENCRYPTION_SECRET — GitHub token 暗号化用 secret
#   NEWS_DATABASE_URL              — news-api 用 Neon DB 接続文字列
#
# 使い方:
#   cp .env.prod.example .env.prod
#   # .env.prod の値を埋める
#   ./scripts/setup-secrets.sh
#
# 別ファイルを使う場合:
#   ENV_FILE=.env.prod ./scripts/setup-secrets.sh
#
# 登録されるシークレット:
#   - lang-war-database-url
#   - lang-war-github-client-id
#   - lang-war-github-client-secret
#   - lang-war-github-redirect-url
#   - lang-war-auth-cookie-secret
#   - lang-war-github-token-encryption-secret
#   - lang-war-news-database-url
set -euo pipefail

# ── プロジェクトルートの .env.prod を読み込む ──────────────────
# スクリプトの場所（scripts/）から一つ上がりプロジェクトルートを特定する
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${ENV_FILE:-${ROOT_DIR}/.env.prod}"

case "${ENV_FILE}" in
  /*) ;;
  *) ENV_FILE="${ROOT_DIR}/${ENV_FILE}" ;;
esac

if [ ! -f "${ENV_FILE}" ]; then
  echo "Error: ${ENV_FILE} が見つかりません"
  echo ""
  echo ".env.prod.example をコピーして値を埋めてください:"
  echo ""
  echo "  cp .env.prod.example .env.prod"
  exit 1
fi

set -a
# shellcheck source=/dev/null
. "${ENV_FILE}"
set +a

# ── 必須変数チェック ──────────────────────────────────────────
missing=""
[ -z "${GCP_PROJECT_ID:-}"       ] && missing="${missing}  GCP_PROJECT_ID\n"
[ -z "${NEON_DATABASE_URL:-}"    ] && missing="${missing}  NEON_DATABASE_URL\n"
[ -z "${GITHUB_CLIENT_ID:-}"     ] && missing="${missing}  GITHUB_CLIENT_ID\n"
[ -z "${GITHUB_CLIENT_SECRET:-}" ] && missing="${missing}  GITHUB_CLIENT_SECRET\n"
[ -z "${GITHUB_REDIRECT_URL:-}"  ] && missing="${missing}  GITHUB_REDIRECT_URL\n"
[ -z "${AUTH_COOKIE_SECRET:-}"   ] && missing="${missing}  AUTH_COOKIE_SECRET\n"
[ -z "${GITHUB_TOKEN_ENCRYPTION_SECRET:-}" ] && missing="${missing}  GITHUB_TOKEN_ENCRYPTION_SECRET\n"
[ -z "${NEWS_DATABASE_URL:-}"              ] && missing="${missing}  NEWS_DATABASE_URL\n"

if [ -n "${missing}" ]; then
  echo "Error: ${ENV_FILE} に以下の変数が設定されていません:"
  printf '%b' "${missing}"
  exit 1
fi

auth_cookie_secret_bytes="$(printf '%s' "${AUTH_COOKIE_SECRET}" | wc -c | tr -d ' ')"
if [ "${auth_cookie_secret_bytes}" -lt 64 ]; then
  echo "Error: AUTH_COOKIE_SECRET は64 bytes以上にしてください。現在: ${auth_cookie_secret_bytes} bytes"
  echo "Chat Service の SECRET_KEY_BASE としても使うため、例: openssl rand -base64 64"
  exit 1
fi

PROJECT_ID="${GCP_PROJECT_ID}"

# ── ヘルパー関数 ──────────────────────────────────────────────

ensure_secret() {
  local name="$1"
  if gcloud secrets describe "${name}" --project="${PROJECT_ID}" &>/dev/null; then
    echo "  (既存のシークレットをスキップ: ${name})"
  else
    gcloud secrets create "${name}" \
      --replication-policy=automatic \
      --project="${PROJECT_ID}"
    echo "  作成: ${name}"
  fi
}

add_secret_version() {
  local name="$1"
  local value="$2"

  if gcloud secrets versions access latest \
      --secret="${name}" \
      --project="${PROJECT_ID}" &>/dev/null; then
    local current_value
    current_value="$(gcloud secrets versions access latest \
      --secret="${name}" \
      --project="${PROJECT_ID}")"

    if [ "${current_value}" = "${value}" ]; then
      echo "  (latest version と同じ値のためスキップ: ${name})"
      return
    fi
  fi

  echo -n "${value}" \
    | gcloud secrets versions add "${name}" \
        --data-file=- \
        --project="${PROJECT_ID}"
  echo "  登録しました。"

  disable_old_versions "${name}"
}

# 最新バージョン以外の enabled version を disable して active version 数を1に保つ
disable_old_versions() {
  local name="$1"
  local latest
  latest="$(gcloud secrets versions list "${name}" \
    --project="${PROJECT_ID}" \
    --filter="state=ENABLED" \
    --format="value(name)" \
    --sort-by="~name" | head -1)"

  gcloud secrets versions list "${name}" \
    --project="${PROJECT_ID}" \
    --filter="state=ENABLED" \
    --format="value(name)" \
  | grep -v "^${latest}$" \
  | while read -r version; do
    gcloud secrets versions disable "${version}" \
      --secret="${name}" \
      --project="${PROJECT_ID}" \
      --quiet
    echo "  古いバージョン ${version} を disabled にしました。"
  done
}

# ── 開始 ─────────────────────────────────────────────────────
echo "===================================="
echo " GCP Secret Manager セットアップ"
echo "===================================="
echo "  Project: ${PROJECT_ID}"
echo "  Env:     ${ENV_FILE}"
echo "===================================="
echo ""

echo "[1/6] lang-war-database-url"
ensure_secret "lang-war-database-url"
add_secret_version "lang-war-database-url" "${NEON_DATABASE_URL}"
echo ""

echo "[2/6] lang-war-github-client-id"
ensure_secret "lang-war-github-client-id"
add_secret_version "lang-war-github-client-id" "${GITHUB_CLIENT_ID}"
echo ""

echo "[3/6] lang-war-github-client-secret"
ensure_secret "lang-war-github-client-secret"
add_secret_version "lang-war-github-client-secret" "${GITHUB_CLIENT_SECRET}"
echo ""

echo "[4/6] lang-war-github-redirect-url"
ensure_secret "lang-war-github-redirect-url"
add_secret_version "lang-war-github-redirect-url" "${GITHUB_REDIRECT_URL}"
echo ""

echo "[5/6] lang-war-auth-cookie-secret"
ensure_secret "lang-war-auth-cookie-secret"
add_secret_version "lang-war-auth-cookie-secret" "${AUTH_COOKIE_SECRET}"
echo ""

echo "[6/7] lang-war-github-token-encryption-secret"
ensure_secret "lang-war-github-token-encryption-secret"
add_secret_version "lang-war-github-token-encryption-secret" "${GITHUB_TOKEN_ENCRYPTION_SECRET}"
echo ""

echo "[7/7] lang-war-news-database-url"
ensure_secret "lang-war-news-database-url"
add_secret_version "lang-war-news-database-url" "${NEWS_DATABASE_URL}"
echo ""

# ── 登録済みシークレット一覧を表示 ───────────────────────────
echo "===================================="
echo " 登録完了"
echo "===================================="
echo ""
gcloud secrets list --project="${PROJECT_ID}" --filter="name:lang-war-" --format="table(name,createTime)"
echo ""
echo "Cloud Run は --set-secrets フラグでこれらを参照します。"
echo "deploy.yml は Secret Manager の latest version を参照します。"
echo ""
