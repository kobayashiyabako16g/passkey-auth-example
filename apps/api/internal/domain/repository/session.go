package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/kvstore"
)

// KV Expire
const Expire = 10

type Session interface {
	Create(ctx context.Context, id string) (*model.Session, error)
	Save(ctx context.Context, session *model.Session) error
	Get(ctx context.Context, id string) (*model.Session, error)
	Delete(ctx context.Context, session *model.Session) error
}

type sessionImpl struct {
	client kvstore.Client
}

func NewSession(client kvstore.Client) Session {
	return &sessionImpl{client}
}

func (s *sessionImpl) Create(ctx context.Context, id string) (*model.Session, error) {
	session := &model.Session{
		ID:        id,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionImpl) Save(ctx context.Context, session *model.Session) error {
	key := s.getKey(session.ID)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	// Calculate TTL in seconds
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	return s.client.Set(ctx, key, string(string(data)), kvstore.SetOptions{Expiration: Expire})
}

func (s *sessionImpl) Get(ctx context.Context, key string) (*model.Session, error) {
	sessionID := s.getKey(key)
	data, err := s.client.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}

	var session model.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}

	return &session, nil
}

func (s *sessionImpl) Delete(ctx context.Context, session *model.Session) error {
	key := s.getKey(session.ID)
	return s.client.Delete(ctx, key)
}

func (s *sessionImpl) getKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
