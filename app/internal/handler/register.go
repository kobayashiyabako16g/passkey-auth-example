package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/server"
)

func BeginRegistration(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
			return
		}

		s.UsersMu.Lock()
		user, exists := s.Users[username]
		if !exists {
			user = &model.User{
				ID:          []byte(username),
				Name:        username,
				DisplayName: username,
				Credentials: []webauthn.Credential{},
			}
			s.Users[username] = user
		}
		s.UsersMu.Unlock()

		options, sessionData, err := s.WebAuthn.BeginRegistration(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
			return
		}

		session := s.GetSession(c)
		if session == nil {
			session = s.SessionStore.Create()
		}
		session.Username = username
		session.RegistrationData = sessionData
		s.SessionStore.Save(session)
		s.SetSessionCookie(c, session)

		c.JSON(http.StatusOK, options)
	}
}

func FinishRegistration(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.GetSession(c)
		if session == nil || session.RegistrationData == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No registration session found"})
			return
		}

		s.UsersMu.RLock()
		user := s.Users[session.Username]
		s.UsersMu.RUnlock()

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		credential, err := s.WebAuthn.FinishRegistration(user, *session.RegistrationData, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish registration"})
			return
		}

		s.UsersMu.Lock()
		user.Credentials = append(user.Credentials, *credential)
		s.UsersMu.Unlock()

		session.RegistrationData = nil
		s.SessionStore.Save(session)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
