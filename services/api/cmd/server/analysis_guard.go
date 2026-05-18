package main

import (
	"context"
	"sync"
	"time"

	analysisapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type analysisGuard struct {
	inner   *analysisapp.UseCase
	session interface {
		FindUserBySessionToken(ctx context.Context, token string, now time.Time) (user.User, bool, error)
	}
	mu    sync.Mutex
	locks map[user.ID]*sync.Mutex
}

func newAnalysisGuard(inner *analysisapp.UseCase, session interface {
	FindUserBySessionToken(ctx context.Context, token string, now time.Time) (user.User, bool, error)
}) *analysisGuard {
	return &analysisGuard{
		inner:   inner,
		session: session,
		locks:   make(map[user.ID]*sync.Mutex),
	}
}

func (g *analysisGuard) Analyze(ctx context.Context, sessionToken string) (analysisapp.AnalysisResult, error) {
	appUser, ok, err := g.session.FindUserBySessionToken(ctx, sessionToken, time.Now())
	if err != nil {
		return analysisapp.AnalysisResult{}, err
	}
	if !ok {
		return analysisapp.AnalysisResult{}, analysisapp.ErrUnauthenticated
	}

	g.mu.Lock()
	l, ok := g.locks[appUser.ID]
	if !ok {
		l = &sync.Mutex{}
		g.locks[appUser.ID] = l
	}
	g.mu.Unlock()

	l.Lock()
	defer func() {
		l.Unlock()
		g.mu.Lock()
		delete(g.locks, appUser.ID)
		g.mu.Unlock()
	}()

	return g.inner.AnalyzeForUser(ctx, appUser, sessionToken, time.Now())
}
