// Package user はアプリユーザー集約と紐づく GitHub アカウント情報を管理する。
package user

import (
	"errors"
	"time"
)

type ID string

type GitHubAccount struct {
	GitHubID  int64
	Username  string
	AvatarURL string
}

type GitHubProfile struct {
	GitHubID  int64
	Username  string
	AvatarURL string
}

type User struct {
	ID            ID
	GitHubAccount GitHubAccount
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewFromGitHub(id ID, profile GitHubProfile, now time.Time) (User, error) {
	if id == "" {
		return User{}, errors.New("user id is required")
	}

	account, err := NewGitHubAccount(profile)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:            id,
		GitHubAccount: account,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (u User) LinkGitHub(profile GitHubProfile, now time.Time) (User, error) {
	account, err := NewGitHubAccount(profile)
	if err != nil {
		return User{}, err
	}

	u.GitHubAccount = account
	u.UpdatedAt = now
	return u, nil
}

func NewGitHubAccount(profile GitHubProfile) (GitHubAccount, error) {
	if profile.GitHubID == 0 {
		return GitHubAccount{}, errors.New("github id is required")
	}
	if profile.Username == "" {
		return GitHubAccount{}, errors.New("github username is required")
	}

	return GitHubAccount(profile), nil
}
