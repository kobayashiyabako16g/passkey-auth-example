package model

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

type Session struct {
	ID                 string
	Username           string
	Authenticated      bool
	RegistrationData   *webauthn.SessionData
	AuthenticationData *webauthn.SessionData
	ExpiresAt          time.Time
}

// SessionStore manages sessions in memory
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*Session),
	}

	// Clean up expired sessions every hour
	go store.cleanup()

	return store
}

func (s *SessionStore) generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (s *SessionStore) Create() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		ID:        s.generateSessionID(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.sessions[session.ID] = session
	return session
}

func (s *SessionStore) Get(sessionID string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}

	return session
}

func (s *SessionStore) Save(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = session
}

func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
}

func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
