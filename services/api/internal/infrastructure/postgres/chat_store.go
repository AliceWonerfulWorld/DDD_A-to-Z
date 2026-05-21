package postgres

import (
	"context"

	chatapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/chat"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

type ChatStore struct {
	db    *gorm.DB
	guild *GuildStore
}

func NewChatStore(db *gorm.DB, guild *GuildStore) *ChatStore {
	return &ChatStore{db: db, guild: guild}
}

func (s *ChatStore) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	return s.guild.FindActiveMembershipByUserID(ctx, userID)
}

func (s *ChatStore) InsertChatToken(ctx context.Context, token chatapp.ChatToken) error {
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO chat_tokens (token_hash, user_id, guild_id, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, token.TokenHash, token.UserID, token.GuildID, token.ExpiresAt, token.CreatedAt).Error
}
