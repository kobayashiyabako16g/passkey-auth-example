package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/repository"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/handler/dto/request"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type Auth interface {
	BeginRegistration(w http.ResponseWriter, r *http.Request)
}

type auth struct {
	sr       repository.Session
	ur       repository.User
	webAuthn *webauthn.WebAuthn
}

func NewAuth(sr repository.Session, ur repository.User) Auth {
	wconfig := &webauthn.Config{
		RPDisplayName: "Passkey Demo",                    // Display Name for your site
		RPID:          "localhost",                       // Generally the domain name for your site
		RPOrigins:     []string{"http://localhost:5173"}, // Vite dev server origin
	}
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}
	return &auth{sr, ur, webAuthn}
}

func (s *auth) setSessionCookie(w http.ResponseWriter, session *model.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Required for cross-origin requests
		Secure:   false,                 // Set to true in production with HTTPS
		Expires:  session.ExpiresAt,
	})
}

func (h *auth) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "begin registration ----------------------")

	var u request.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logger.Info(ctx, "can't decode user data", logger.WithError(err))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	if u.Username == "" {
		logger.Info(ctx, "username is empty")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// ユーザー登録(確認)
	exists, err := h.ur.ExistsByUsername(ctx, u.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		logger.Info(ctx, fmt.Sprintf("Exists User name: %s", u.Username))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	var user model.User
	if err = user.GenerateID(); err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Name = u.Username
	user.DisplayName = u.Username

	// チャレンジ生成
	options, sessionData, err := h.webAuthn.BeginRegistration(user)
	if err != nil {
		logger.Error(ctx, "Error beginning registration", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, err := h.sr.Create(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to create session", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Username = u.Username
	session.RegistrationData = sessionData

	// Store に保存
	err = h.sr.Save(ctx, session)
	if err != nil {
		logger.Error(ctx, "Failed to store challenge", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.setSessionCookie(w, session)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(options); err != nil {
		logger.Error(ctx, "Failed to write response", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// TODO: Environment Value
var jwtSecret = []byte("your-very-secret-key")

func GenerateRegisterJWT(userID string, challenge string) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userID,
		"challenge": challenge,
		"exp":       time.Now().Add(5 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
