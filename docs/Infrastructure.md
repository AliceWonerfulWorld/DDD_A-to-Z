# インフラ運用ガイド — Lang War

## 全体構成

```
GitHub Actions (CI/CD)
  │
  ├── PR 時: format / lint / test / build（自動）
  │
  └── main マージ時: DB migration → Docker build → Cloud Run deploy（自動）

開発者のローカル（初回セットアップ時のみ手動）
  │
  ├── tofu apply       → GCP リソースを作成・変更
  └── setup-secrets.sh → Secret Manager に値を登録
```

---

## ファイル構成

```
infra/
└── terraform/
    ├── main.tf               # OpenTofu 自体の設定（state 保存先など）
    ├── variables.tf          # 外から渡す変数（project_id 等）
    ├── apis.tf               # GCP API の有効化
    ├── wif.tf                # GitHub Actions → GCP 認証（Workload Identity Federation）
    ├── artifact_registry.tf  # Docker イメージの保存場所
    ├── secrets.tf            # Secret Manager のシークレット定義（箱だけ）
    ├── cloud_run.tf          # Cloud Run サービスの定義
    └── outputs.tf            # tofu apply 後に表示される値

scripts/
├── atlas.sh          # Atlas CLI のラッパー（.env を自動読み込み）
└── setup-secrets.sh  # Secret Manager に値を登録するスクリプト

.github/workflows/
├── ci-api.yml    # Go API の format / lint / test / build
├── ci-web.yml    # フロントエンドの format / lint / test / build
├── ci-db.yml     # migration の validate / lint
├── ci-proto.yml  # buf lint / build / generate / breaking change 検出
├── ci-infra.yml  # flake.nix の検証
└── deploy.yml    # main マージ時の自動デプロイ
```

---

## 初回セットアップ（一度だけ実行）

### 前提

- `gcloud auth login` でオーナー権限のアカウントで認証済み
- プロジェクトルートの `.env.prod` に本番 Secret 登録用の変数が設定済み（`.env.prod.example` を参照）

```bash
cp .env.prod.example .env.prod
# エディタで実際の値を埋める
```

`.env.prod` に設定する変数：

| 変数 | 説明 |
|---|---|
| `GCP_PROJECT_ID` | GCP プロジェクト ID |
| `NEON_DATABASE_URL` | Neon DB の接続文字列 |
| `GITHUB_CLIENT_ID` | GitHub OAuth App の Client ID |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth App の Client Secret |
| `GITHUB_REDIRECT_URL` | GitHub OAuth callback URL |
| `AUTH_COOKIE_SECRET` | OAuth state cookie 署名用 secret。Chat Service の `SECRET_KEY_BASE` としても使うため64 bytes以上 |
| `GITHUB_TOKEN_ENCRYPTION_SECRET` | GitHub token 暗号化用 secret |

Chat Service の `SECRET_KEY_BASE` は Secret Manager の数を増やさないため、既存の `AUTH_COOKIE_SECRET` を流用する。Phoenix の本番起動要件に合わせて、`AUTH_COOKIE_SECRET` は64 bytes以上にする。
`CHAT_ALLOWED_ORIGIN` は秘密情報ではないため、Terraform 変数 `chat_allowed_origin` と deploy workflow の通常環境変数で管理する。

### 手順

**1. tfstate 用の GCS バケットを作成**

OpenTofu の状態ファイル置き場。1回だけ手動で作る。

```bash
gcloud storage buckets create gs://PROJECT_ID-tfstate \
  --location=asia-northeast1 \
  --project=PROJECT_ID
```

**2. OpenTofu でGCP リソースを作成**

```bash
cd infra/terraform

tofu init -backend-config="bucket=PROJECT_ID-tfstate"
tofu plan -var="project_id=PROJECT_ID" -var="github_repo=jyogi-web/DDD_A-to-Z"
tofu apply -var="project_id=PROJECT_ID" -var="github_repo=jyogi-web/DDD_A-to-Z"
```

`tofu apply` 完了後、以下の値が表示される。

```
Outputs:

workload_identity_provider = "projects/123.../providers/github-actions-provider"
service_account_email      = "lang-war-deploy@PROJECT_ID.iam.gserviceaccount.com"
cloud_run_url              = "https://lang-war-api-xxxx-an.a.run.app"
chat_service_url           = "https://chat-service-xxxx-an.a.run.app"
artifact_registry_repo     = "asia-northeast1-docker.pkg.dev/PROJECT_ID/lang-war/api"
chat_artifact_registry_repo = "asia-northeast1-docker.pkg.dev/PROJECT_ID/lang-war/chat-service"
```

**3. GitHub にシークレットと変数を登録**

`Settings > Secrets and variables > Actions` に以下を登録する。

| 種別 | 名前 | 値 |
|---|---|---|
| Variables | `GCP_PROJECT_ID` | GCP プロジェクト ID |
| Secrets | `GCP_WORKLOAD_IDENTITY_PROVIDER` | `tofu output workload_identity_provider` の値 |
| Secrets | `GCP_SERVICE_ACCOUNT` | `tofu output service_account_email` の値 |
| Secrets | `NEON_DATABASE_URL` | Neon DB の接続文字列 |

**4. Secret Manager に値を登録**

```bash
./scripts/setup-secrets.sh
```

`.env.prod` の値を読んで Secret Manager に登録する。以降、Cloud Run 起動時に自動で注入される。

---

## 普段の開発フロー

```
feature ブランチで開発
  └── git push → PR を作成
        └── CI が自動実行（format / lint / test / build / migrate validate）
              └── 全 CI グリーン → main にマージ
                    └── deploy.yml が自動実行
                          ├── atlas migrate apply（Neon DB に migration 適用）
                          ├── docker build & push（Artifact Registry）
                          └── gcloud run deploy（Cloud Run のイメージを差し替え）
```

main への直接 push は禁止。必ず PR を経由する。

---

## インフラを変更したいとき

コードの変更（API の機能追加など）は上記の普段の開発フローで完結する。
以下の場合だけ、ローカルで `tofu apply` を追加で実行する。

| 変更内容 | 対応するファイル | 実行コマンド |
|---|---|---|
| Cloud Run のメモリ・CPU・スケール設定を変えたい | `cloud_run.tf` を編集 | `tofu apply` |
| 新しいシークレットを追加したい | `secrets.tf` を編集 → `setup-secrets.sh` で値を登録 | `tofu apply` → `./scripts/setup-secrets.sh` |
| 新しい GCP API を使いたい | `apis.tf` に追加 | `tofu apply` |
| 新しい Cloud Run サービスを追加したい | `cloud_run.tf` に追加（将来は module 化） | `tofu apply` |

`tofu apply` の前に必ず `tofu plan` で差分を確認すること。

```bash
cd infra/terraform
tofu plan -var="project_id=PROJECT_ID" -var="github_repo=jyogi-web/DDD_A-to-Z"
```

---

## DB migration の運用

migration ファイルは `db/migrations/` に置き、Atlas で管理する。

### 新しい migration を作る

```bash
pnpm db:migrate:diff <migration名>
# 例: pnpm db:migrate:diff add_season_table
```

### ローカル DB に適用して動作確認

```bash
pnpm db:up              # Docker で PostgreSQL を起動
pnpm db:migrate:apply   # ローカル DB に適用
```

### PR を出す

CI の `ci-db.yml` が自動で以下を実行する。

- `migrate validate` — migration ファイルの整合性チェック
- `migrate lint` — 破壊的変更（DROP TABLE 等）の検出

CI が通ったら PR をマージ。`deploy.yml` が Neon DB（本番）に migration を適用する。

---

## 障害対応

### Cloud Run が起動しない

```bash
# ログを確認
gcloud run services logs read lang-war-api --region=asia-northeast1 --limit=50
```

よくある原因：

- Secret Manager のシークレットに値が登録されていない → `./scripts/setup-secrets.sh` を再実行
- Neon DB への接続失敗 → `NEON_DATABASE_URL` の値を確認
- migration が失敗している → `deploy.yml` の migration ステップのログを確認

### migration をロールバックしたい

Atlas は down migration を自動生成しない。手動で巻き戻し用の migration ファイルを作成して PR を出す。

```bash
pnpm db:migrate:diff revert_<migration名>
```

---

## 関連ドキュメント

- `docs/Architecture.md` — システム全体のアーキテクチャ
- `docs/BackendStructure.md` — Go API の層構造
- `infra/terraform/` — 各 `.tf` ファイルにコメントで詳細を記載
