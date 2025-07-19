package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/server"
)

func Logout(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.GetSession(c)
		if session != nil {
			s.SessionStore.Delete(session.ID)
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			MaxAge:   -1,
		})

		c.JSON(http.StatusOK, gin.H{"status": "logged out"})
	}
}
