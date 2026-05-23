package mypage

import (
	"context"

	domainprofile "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/profile"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// ProfileRepository is the interface for accessing profile domain models.
type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID user.ID) (domainprofile.Profile, bool, error)
}

type profileReader struct {
	repo ProfileRepository
}

// NewProfileReader creates a ProfileReader that wraps a ProfileRepository.
func NewProfileReader(repo ProfileRepository) ProfileReader {
	return &profileReader{repo: repo}
}

func (r *profileReader) GetProfile(ctx context.Context, userID user.ID) (*ProfileInfo, error) {
	prof, exists, err := r.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	return &ProfileInfo{
		DisplayName: prof.DisplayName,
		AvatarURL:   prof.AvatarURL,
	}, nil
}
