package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/server"
)

func Info(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := s.GetSession(c)
		if session == nil || !session.Authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to your dashboard!",
			"user":    session.Username,
			"time":    time.Now(),
		})
	}
}
