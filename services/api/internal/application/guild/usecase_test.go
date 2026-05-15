package guild

import (
	"context"
	"fmt"
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type testRepository struct {
	guilds []guilddomain.Guild
}

func (r testRepository) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	return r.guilds, nil
}

func TestUseCaseListGuilds(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	expected := []guilddomain.Guild{{
		ID:          "guild_typescript",
		Slug:        "typescript",
		Name:        "TypeScript",
		Description: "型の力で支えるギルド。",
		Icon:        "TS",
		Color:       "#3178c6",
		SortOrder:   5,
		CreatedAt:   now,
		UpdatedAt:   now,
	}}
	usecase := NewUseCase(testRepository{guilds: expected})

	guilds, err := usecase.ListGuilds(context.Background())
	if err != nil {
		t.Fatalf("ListGuilds() がエラーを返しました: %v", err)
	}
	if len(guilds) != 1 {
		t.Fatalf("guilds length = %d, 期待値 1", len(guilds))
	}
	if guilds[0].Slug != "typescript" {
		t.Fatalf("Slug = %q, 期待値 typescript", guilds[0].Slug)
	}
}

func TestNewUseCasePanicsWithoutRepository(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("NewUseCase() panic = nil, 期待値 panic")
		}
		if message := fmt.Sprint(recovered); message != "guild repository is required" {
			t.Fatalf("NewUseCase() panic = %q, 期待値 guild repository is required", message)
		}
	}()

	_ = NewUseCase(nil)
}
