package server

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
)

type Server struct {
	WebAuthn     *webauthn.WebAuthn
	SessionStore *model.SessionStore
	Users        map[string]*model.User
	UsersMu      sync.RWMutex
}

func (s *Server) GetSession(c *gin.Context) *model.Session {
	cookie, err := c.Request.Cookie("session")
	if err != nil {
		return nil
	}
	return s.SessionStore.Get(cookie.Value)
}

func (s *Server) SetSessionCookie(c *gin.Context, session *model.Session) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
		Expires:  session.ExpiresAt,
	})
}

func NewServer() (*Server, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "My Passkey App",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	})
	if err != nil {
		return nil, err
	}

	return &Server{
		WebAuthn:     webAuthn,
		SessionStore: model.NewSessionStore(),
		Users:        make(map[string]*model.User),
	}, nil
}
