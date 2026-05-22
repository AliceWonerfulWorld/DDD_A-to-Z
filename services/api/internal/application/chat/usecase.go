package chat

import (
	"context"
	"errors"
	"strings"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden: not a member of this guild")
)

const chatTokenTTL = 5 * time.Minute

type ChatToken struct {
	Token     string
	TokenHash string
	UserID    string
	GuildID   string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UseCase struct {
	current CurrentUserRepository
	repo    Repository
	tokens  TokenGenerator
	hasher  Hasher
	now     func() time.Time
}

func NewUseCase(current CurrentUserRepository, repo Repository, tokens TokenGenerator, hasher Hasher) *UseCase {
	if current == nil {
		panic("current user repository is required")
	}
	if repo == nil {
		panic("chat repository is required")
	}
	if tokens == nil {
		panic("token generator is required")
	}
	if hasher == nil {
		panic("hasher is required")
	}

	return &UseCase{
		current: current,
		repo:    repo,
		tokens:  tokens,
		hasher:  hasher,
		now:     time.Now,
	}
}

func (u *UseCase) IssueGuildChatToken(ctx context.Context, sessionToken string, guildID guilddomain.ID) (ChatToken, error) {
	normalizedSessionToken := strings.TrimSpace(sessionToken)
	if normalizedSessionToken == "" {
		return ChatToken{}, ErrUnauthenticated
	}
	if strings.TrimSpace(string(guildID)) == "" {
		return ChatToken{}, ErrForbidden
	}

	now := u.now()
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, normalizedSessionToken, now)
	if err != nil {
		return ChatToken{}, err
	}
	if !ok {
		return ChatToken{}, ErrUnauthenticated
	}

	_, ok, err = u.repo.FindMembershipByUserAndGuild(ctx, appUser.ID, guildID)
	if err != nil {
		return ChatToken{}, err
	}
	if !ok {
		return ChatToken{}, ErrForbidden
	}

	rawToken, err := u.tokens.NewToken()
	if err != nil {
		return ChatToken{}, err
	}

	token := ChatToken{
		Token:     rawToken,
		TokenHash: u.hasher.Hash(rawToken),
		UserID:    string(appUser.ID),
		GuildID:   string(guildID),
		ExpiresAt: now.Add(chatTokenTTL),
		CreatedAt: now,
	}

	if err := u.repo.InsertChatToken(ctx, token); err != nil {
		return ChatToken{}, err
	}

	return token, nil
}
