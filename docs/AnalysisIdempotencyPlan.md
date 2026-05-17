# 解析の重複換算防止 実装計画

## 課題

現在 `Analyze()` は毎回 `now - 30日` を起点にコミット/PRを取得し、その都度 CP を `Earn()` する。同じ期間を複数回実行すると同じ contribution が重複換算されてしまう。

## 解決策

各ユーザーの「最終換算日時 (`last_analyzed_at`)」を記録し、次回解析時はその日時以降の contribution のみを対象にする。

---

## 1. DB Schema 変更

### 対象テーブル

`point_accounts` にカラム追加:

```sql
ALTER TABLE point_accounts
ADD COLUMN last_analyzed_at TIMESTAMPTZ;
```

- 初回は NULL（= 未解析）
- `db/schema/points.sql` にも追記

### 新規マイグレーションファイル

`db/migrations/20260517000000_add_last_analyzed_at.sql` を作成。

---

## 2. Port 変更

### `application/repositoryanalysis/ports.go`

`CPBalanceProvider` にメソッド追加（改名はしない）:

```go
type CPBalanceProvider interface {
    GetBalance(ctx context.Context, userID user.ID) (int64, error)
    GetLastAnalyzedAt(ctx context.Context, userID user.ID) (*time.Time, error)
    UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, at time.Time) error
}
```

※ `PointType` はこの interface では扱わず、`cpManager` アダプター側で `PointTypeCP` を注入する。

### `application/contributionpoint/ports.go`

`LedgerRepository` にメソッド追加（既存の `GetBalance` と同様に `pointType` を受け取る）:

```go
type LedgerRepository interface {
    Record(ctx context.Context, entry contributionpointdomain.LedgerEntry) (contributionpointdomain.LedgerEntry, error)
    GetBalance(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (int64, error)
    GetLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (*time.Time, error)
    UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType, at time.Time) error
}
```

---

## 3. UseCase 変更

### `application/repositoryanalysis/usecase.go`

```go
// since 計算を last_analyzed_at ベースに変更:
since := now.Add(analysisPeriod)
lastAnalyzedAt, err := u.cpBalance.GetLastAnalyzedAt(ctx, appUser.ID)
if err != nil {
    return AnalysisResult{}, err
}
if lastAnalyzedAt != nil && lastAnalyzedAt.After(since) {
    since = *lastAnalyzedAt
}

// Earn 後に last_analyzed_at を更新 (totalCP==0 でも更新):
if totalCP > 0 {
    if err := u.cp.Earn(ctx, appUser.ID, totalCP, ...); err != nil {
        return AnalysisResult{}, err
    }
}
if err := u.cpBalance.UpdateLastAnalyzedAt(ctx, appUser.ID, now); err != nil {
    return AnalysisResult{}, err
}
```

### `application/contributionpoint/usecase.go`

パススルーメソッド追加（`GetBalance` と同様に `pointType` を受け取る）:

```go
func (u *UseCase) GetLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (*time.Time, error) {
    return u.ledger.GetLastAnalyzedAt(ctx, userID, pointType)
}

func (u *UseCase) UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType, at time.Time) error {
    return u.ledger.UpdateLastAnalyzedAt(ctx, userID, pointType, at)
}
```

---

## 4. Infrastructure 実装

### `infrastructure/postgres/contributionpoint_store.go`

`GetLastAnalyzedAt` / `UpdateLastAnalyzedAt` を実装:

```go
func (s *ContributionPointStore) GetLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType) (*time.Time, error) {
    var record struct {
        LastAnalyzedAt *time.Time `gorm:"column:last_analyzed_at"`
    }
    result := s.db.WithContext(ctx).
        Select("last_analyzed_at").
        Where("user_id = ? AND point_type = ?", userID, pointType).
        Take(&record)
    if result.Error != nil {
        return nil, result.Error
    }
    return record.LastAnalyzedAt, nil
}

func (s *ContributionPointStore) UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, pointType contributionpointdomain.PointType, at time.Time) error {
    return s.db.WithContext(ctx).
        Model(&pointAccountRecord{}).
        Where("user_id = ? AND point_type = ?", userID, pointType).
        Update("last_analyzed_at", at).Error
}
```

`pointAccountRecord` にフィールド追加:

```go
type pointAccountRecord struct {
    UserID         string                            `gorm:"column:user_id"`
    PointType      contributionpointdomain.PointType `gorm:"column:point_type"`
    Balance        int64                             `gorm:"column:balance"`
    LastAnalyzedAt *time.Time                        `gorm:"column:last_analyzed_at"`
}
```

---

## 5. DI / main.go

### `cmd/server/main.go`

`cpManager` にメソッド追加（`GetBalance` と同様に内側で `PointTypeCP` を注入）:

```go
func (m *cpManager) GetLastAnalyzedAt(ctx context.Context, userID user.ID) (*time.Time, error) {
    return m.inner.GetLastAnalyzedAt(ctx, userID, contributionpointdomain.PointTypeCP)
}

func (m *cpManager) UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, at time.Time) error {
    return m.inner.UpdateLastAnalyzedAt(ctx, userID, contributionpointdomain.PointTypeCP, at)
}
```

---

## 影響範囲まとめ

| ファイル | 変更内容 |
|---|---|
| `db/schema/points.sql` | `point_accounts` に `last_analyzed_at` 追加 |
| `db/migrations/20260517000000_add_last_analyzed_at.sql` | 新規マイグレーションファイル |
| `application/repositoryanalysis/ports.go` | `CPBalanceProvider` に `GetLastAnalyzedAt`/`UpdateLastAnalyzedAt` 追加 |
| `application/repositoryanalysis/usecase.go` | `since` 計算変更、日時更新ロジック追加 |
| `application/contributionpoint/ports.go` | `LedgerRepository` に `pointType` 付きメソッド追加 |
| `application/contributionpoint/usecase.go` | パススルーメソッド追加 |
| `infrastructure/postgres/contributionpoint_store.go` | `GetLastAnalyzedAt`/`UpdateLastAnalyzedAt` 実装、`pointAccountRecord` にフィールド追加 |
| `cmd/server/main.go` | `cpManager` に `GetLastAnalyzedAt`/`UpdateLastAnalyzedAt` 追加（内側で `PointTypeCP` 注入） |

## 変更不要

- `interfaces/http/analysis_controller.go`（レスポンス形式は変わらない）
- `apps/web/`（フロントエンドは影響なし）
- `domain/`（ドメインモデルに変更なし）

## 動作イメージ

### 初回解析
1. `last_analyzed_at` が NULL → `since = now - 30日`
2. コミット/PRを取得・換算
3. CP を `Earn()`
4. `point_accounts.last_analyzed_at = now` に更新（`point_type = 'CP'`）

### 2回目以降の解析
1. `last_analyzed_at` が前回日時 → `since = last_analyzed_at`
2. その日時以降のコミット/PRのみ取得・換算
3. 重複なく新規分だけ CP 加算
4. `point_accounts.last_analyzed_at = now` に更新

### 長期未使用ユーザー
- NULL のままなら `now - 30日` の30日ウィンドウを維持
- 古い `last_analyzed_at` が残っていても `since = max(last_analyzed_at, now - 30日)` で30日を超えない
