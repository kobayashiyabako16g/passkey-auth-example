package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/server"
)

func BeginLogin(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
			return
		}

		s.UsersMu.RLock()
		user := s.Users[username]
		s.UsersMu.RUnlock()

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		options, sessionData, err := s.WebAuthn.BeginLogin(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin login"})
			return
		}

		session := s.GetSession(c)
		if session == nil {
			session = s.SessionStore.Create()
		}
		session.Username = username
		session.AuthenticationData = sessionData
		s.SessionStore.Save(session)
		s.SetSessionCookie(c, session)

		c.JSON(http.StatusOK, options)
	}
}

func FinishLogin(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.GetSession(c)
		if session == nil || session.AuthenticationData == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No authentication session found"})
			return
		}

		s.UsersMu.RLock()
		user := s.Users[session.Username]
		s.UsersMu.RUnlock()

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		_, err := s.WebAuthn.FinishLogin(user, *session.AuthenticationData, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish login"})
			return
		}

		session.Authenticated = true
		session.AuthenticationData = nil
		s.SessionStore.Save(session)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
