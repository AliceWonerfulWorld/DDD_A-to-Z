# ============================================================
# outputs.tf — terraform apply 後に表示される値
# ============================================================
#
# output は terraform apply 完了後にターミナルに表示される値。
# 「作成したリソースの情報をどこかに貼り付けたい」時に便利。
#
# terraform apply 後に以下のように表示される:
#
#   Outputs:
#
#   workload_identity_provider = "projects/123456/locations/global/..."
#   service_account_email      = "lang-war-deploy@my-project.iam.gserviceaccount.com"
#   cloud_run_url              = "https://lang-war-api-xxxx-an.a.run.app"
#   artifact_registry_repo     = "asia-northeast1-docker.pkg.dev/my-project/lang-war/api"
#
# tofu output コマンドで後から再表示もできる:
#   tofu output workload_identity_provider
# ============================================================

# GitHub Secrets に登録する値
output "workload_identity_provider" {
  description = "GitHub Secrets の GCP_WORKLOAD_IDENTITY_PROVIDER に設定する"
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "service_account_email" {
  description = "GitHub Secrets の GCP_SERVICE_ACCOUNT に設定する"
  value       = google_service_account.deploy.email
}

# 確認用
output "cloud_run_url" {
  description = "デプロイ後の API エンドポイント URL"
  value       = google_cloud_run_v2_service.api.uri
}

output "chat_service_url" {
  description = "デプロイ後の Chat Service エンドポイント URL"
  value       = google_cloud_run_v2_service.chat_service.uri
}

output "news_api_url" {
  description = "デプロイ後の News API エンドポイント URL"
  value       = google_cloud_run_v2_service.news_api.uri
}

output "artifact_registry_repo" {
  description = "docker push の宛先（タグなし）"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/lang-war/api"
}

output "chat_artifact_registry_repo" {
  description = "chat-service docker push の宛先（タグなし）"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/lang-war/chat-service"
}
