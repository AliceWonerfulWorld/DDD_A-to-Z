# ============================================================
# cloud_run.tf — Cloud Run サービスの定義
# ============================================================
#
# Cloud Run は Docker コンテナをサーバーレスで動かす GCP サービス。
# サーバーの管理不要、リクエストがなければインスタンス数ゼロになりコストゼロ。
#
# ── Terraform と deploy.yml の役割分担 ───────────────────
# Terraform: サービスの「設定」を管理（メモリ、CPU、スケール、環境変数など）
# deploy.yml: サービスの「イメージ」だけを差し替える（gcloud run deploy）
#
# イメージの差し替えは Terraform 管轄外なので、
# deploy.yml で gcloud run deploy を実行しても Terraform の設定は上書きされない。
#
# ── スケール設定について ──────────────────────────────────
# min_instance_count = 0:
#   リクエストがなければインスタンスを0にする（コスト削減）
#   MVPフェーズではアクセスが少ないため0からスタート
#   コールドスタート（初回リクエストが遅くなる）が気になったら 1 に上げる
#
# max_instance_count = 10:
#   同時に最大10インスタンスまでスケールアウト
#   MVP規模では十分。トラフィックが増えたら上げる
#
# ── リソース設定について ──────────────────────────────────
# memory = "512Mi": MVP では十分。Go は省メモリなので余裕がある
# cpu = "1":        1 vCPU。負荷が上がったら "2" に変更する
# timeout = "30s":  リクエストの最大処理時間。GitHub API 呼び出しを含むので30秒
# ============================================================

locals {
  # 初回 tofu apply 時はまだ Artifact Registry にイメージがないため placeholder を使う。
  # main に push した時に deploy.yml が実際のイメージに差し替える。
  image = "us-docker.pkg.dev/cloudrun/container/hello:latest"

  # 環境変数の定義。キーがアプリ内での環境変数名、値が Secret Manager のシークレット名。
  api_secret_envs = {
    DATABASE_URL                   = "lang-war-database-url"
    GITHUB_CLIENT_ID               = "lang-war-github-client-id"
    GITHUB_CLIENT_SECRET           = "lang-war-github-client-secret"
    GITHUB_REDIRECT_URL            = "lang-war-github-redirect-url"
    AUTH_COOKIE_SECRET             = "lang-war-auth-cookie-secret"
    GITHUB_TOKEN_ENCRYPTION_SECRET = "lang-war-github-token-encryption-secret"
  }

  api_plain_envs = {
    AUTH_COOKIE_SECURE = "true"
    FRONTEND_URL       = "https://ddd-a-to-z-web.vercel.app/"
  }

  chat_secret_envs = {
    DATABASE_URL    = "lang-war-database-url"
    SECRET_KEY_BASE = "lang-war-auth-cookie-secret"
  }

  chat_plain_envs = {
    ALLOWED_ORIGIN = var.chat_allowed_origin
  }
}

# Go API サービス
resource "google_cloud_run_v2_service" "api" {
  project  = var.project_id
  name     = "lang-war-api"
  location = var.region

  # false にしないと terraform destroy でサービスを削除できない
  # 誤削除防止のためデフォルトは true だが、開発中は false にしておく
  deletion_protection = false

  template {
    service_account = google_service_account.cloud_run.email

    scaling {
      min_instance_count = 0
      max_instance_count = 10
    }

    containers {
      image = local.image

      ports {
        container_port = 8080 # Go サーバーが listen するポート
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      # 環境変数は Secret Manager から参照する
      # アプリは os.Getenv("DATABASE_URL") で普通に読める
      # 値は起動時に Secret Manager から自動取得されるため、
      # コードや設定ファイルに平文で書く必要がない

      dynamic "env" {
        for_each = local.api_secret_envs

        content {
          name = env.key

          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      dynamic "env" {
        for_each = local.api_plain_envs

        content {
          name  = env.key
          value = env.value
        }
      }
    }

    timeout = "30s"
  }

  # Secret Manager の箱と Artifact Registry が作られてから Cloud Run を作る
  depends_on = [
    google_artifact_registry_repository.api,
    google_secret_manager_secret.secrets,
  ]

  # GCP が scaling ブロックに manual_instance_count を自動付与するため
  # tofu plan のたびに差分として検出される。実態は変わらないので無視する。
  lifecycle {
    ignore_changes = [
      # deploy.yml が実イメージを差し替えるため、Terraform は placeholder との差分を無視する。
      template[0].containers[0].image,
      # GCP がトップレベルに scaling ブロックを自動付与するため
      # plan のたびに差分として検出される。実態は変わらないので無視する。
      scaling,
    ]
  }
}

# ── パブリックアクセス設定 ────────────────────────────────
# Cloud Run はデフォルトで認証が必要。
# Lang War の API は外部からアクセスできる必要があるため、
# allUsers（誰でも）に roles/run.invoker（呼び出し権限）を付与する。
resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Elixir チャットサービス
resource "google_cloud_run_v2_service" "chat_service" {
  project  = var.project_id
  name     = "chat-service"
  location = var.region

  deletion_protection = false

  template {
    service_account = google_service_account.cloud_run.email

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = local.image

      ports {
        container_port = 4000 # Phoenix が listen するポート
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      dynamic "env" {
        for_each = local.chat_secret_envs

        content {
          name = env.key

          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      dynamic "env" {
        for_each = local.chat_plain_envs

        content {
          name  = env.key
          value = env.value
        }
      }
    }

    # WebSocket 接続も Cloud Run では request として扱われるため、最大値まで伸ばす。
    timeout = "3600s"
  }

  depends_on = [
    google_artifact_registry_repository.api,
    google_secret_manager_secret.secrets,
  ]

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
      scaling,
    ]
  }
}

resource "google_cloud_run_v2_service_iam_member" "chat_service_public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.chat_service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
