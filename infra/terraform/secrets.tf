# ============================================================
# secrets.tf — Secret Manager のシークレット定義
# ============================================================
#
# Secret Manager は DB パスワードや API キーなど機密情報を安全に管理する GCP サービス。
# 環境変数に平文で書いたり git に混入させたりしない。
#
# ── Terraform が管理するのは「箱」だけ ───────────────────
# シークレットには「箱（名前と設定）」と「中身（実際の値）」の2層がある。
#
#   箱: Terraform で管理 ← このファイルで定義
#   中身: scripts/setup-secrets.sh で登録
#
# 中身を Terraform で管理しない理由:
#   - terraform.tfstate に平文で記録されてしまう
#   - tfstate を GCS に置いても漏洩リスクが残る
#   - 機密値はコードや状態ファイルから完全に分離したい
#
# ── Cloud Run との連携 ────────────────────────────────────
# cloud_run.tf で env ブロックに secret_key_ref を設定することで、
# Cloud Run 起動時に Secret Manager から自動的に値を取得して
# 環境変数として注入される。アプリは普通の環境変数として読めばよい。
# ============================================================

locals {
  secret_names = [
    "lang-war-database-url",                   # Neon DB 接続文字列
    "lang-war-github-client-id",               # GitHub OAuth App の Client ID
    "lang-war-github-client-secret",           # GitHub OAuth App の Client Secret
    "lang-war-github-redirect-url",            # GitHub OAuth callback URL
    "lang-war-auth-cookie-secret",             # OAuth state cookie 署名用 secret
    "lang-war-github-token-encryption-secret", # GitHub token 暗号化用 secret
  ]
}

resource "google_secret_manager_secret" "secrets" {
  for_each = toset(local.secret_names)

  project   = var.project_id
  secret_id = each.value

  # replication auto = GCP が自動でリージョンを選んで冗長化する
  # 特定リージョンに固定したい場合は user_managed を使う
  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]

  lifecycle {
    prevent_destroy = true
  }
}
