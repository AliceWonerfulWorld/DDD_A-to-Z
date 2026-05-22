# ============================================================
# variables.tf — 外から渡す設定値の定義
# ============================================================
#
# Terraform の変数定義ファイル。
# ハードコードしたくない値（プロジェクト ID など環境ごとに変わるもの）を
# 変数として定義しておき、terraform apply 時に外から渡す。
#
# 渡し方は3通りある（優先度順）：
#   1. コマンドライン引数:  terraform apply -var="project_id=my-project"
#   2. terraform.tfvars:   project_id = "my-project" と書いたファイル
#   3. 環境変数:           TF_VAR_project_id=my-project
#
# default が設定されていない変数は必須。apply 時に入力を求められる。
# ============================================================

variable "project_id" {
  description = "GCP プロジェクト ID（例: lang-war-123456）"
  type        = string
  # default なし = 必須。terraform apply 時に -var で渡すか、対話入力する
}

variable "region" {
  description = "リソースを作成する GCP リージョン"
  type        = string
  default     = "asia-northeast1" # 東京リージョン。ユーザーが日本にいるため
}

variable "github_repo" {
  description = "デプロイ元の GitHub リポジトリ（例: jyogi-web/DDD_A-to-Z）"
  type        = string
  # Workload Identity Federation でこのリポジトリからのみ GCP 操作を許可するために使う
}

variable "chat_allowed_origin" {
  description = "Chat Service の WebSocket 接続を許可する frontend origin"
  type        = string
  default     = "https://ddd-a-to-z-web.vercel.app"
}
