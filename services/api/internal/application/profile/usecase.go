// プロフィール初期設定フローの管理
package profile

import (
	"context"
	"errors"
	"time"

	domainprofile "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/profile"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// エラーパッケージ(未ログイン、プロフィールが既に存在、無効な表示名)
var (
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrProfileExists      = errors.New("profile already completed")
	ErrInvalidDisplayName = errors.New("invalid display name")
)

// UseCase本体(現在ユーザー取得、プロフィール保存)
type UseCase struct {
	current  CurrentUserRepository
	profiles ProfileRepository
	now      func() time.Time
}

// NewUseCase creates a new profile use case.
func NewUseCase(current CurrentUserRepository, profiles ProfileRepository) *UseCase {
	return &UseCase{
		current:  current,
		profiles: profiles,
		now:      time.Now,
	}
}

// 入力DTO(セッショントークン、表示名)
type CompleteInitialProfileInput struct {
	SessionToken string
	DisplayName  string
}

// 認可済みユーザーのプロフィールを初期設定する(プロフィールが既に存在する場合はErrProfileExistsを返す)
func (u *UseCase) CompleteInitialProfile(ctx context.Context, input CompleteInitialProfileInput) error {
	if input.SessionToken == "" {
		return ErrUnauthenticated
	}

	// セッショントークンからユーザーを取得
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, input.SessionToken, u.now())
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnauthenticated
	}

	// ユーザーIDからプロフィールを取得
	_, exists, err := u.profiles.FindByUserID(ctx, appUser.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrProfileExists
	}

	// ドメインエンティティを作成(表示名を検証)
	p, err := domainprofile.New(appUser.ID, input.DisplayName, u.now())
	if err != nil {
		return errors.Join(ErrInvalidDisplayName, err)
	}

	// 4. Persist
	return u.profiles.Save(ctx, p)
}

// FindUser resolves a session token to the authenticated user.
func (u *UseCase) FindUser(ctx context.Context, sessionToken string) (user.User, bool, error) {
	return u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
}

// GetProfile returns the profile for the given user ID.
func (u *UseCase) GetProfile(ctx context.Context, userID user.ID) (domainprofile.Profile, bool, error) {
	return u.profiles.FindByUserID(ctx, userID)
}
