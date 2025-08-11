package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/repository"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/handler/dto/request"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type Auth interface {
	BeginRegistration(w http.ResponseWriter, r *http.Request)
	FinishRegistration(w http.ResponseWriter, r *http.Request)
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

func (s *auth) setSessionCookie(w *http.ResponseWriter, session *model.Session) {
	http.SetCookie(*w, &http.Cookie{
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

	var req request.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info(ctx, "can't decode user data", logger.WithError(err))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		logger.Info(ctx, "username is empty")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// ユーザー確認
	exists, err := h.ur.ExistsByUsername(ctx, req.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		logger.Info(ctx, fmt.Sprintf("Exists User name: %s", req.Username))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// ユーザ作成
	var user model.User
	if err = user.GenerateID(); err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Name = req.Username
	user.DisplayName = req.Username

	// チャレンジ生成
	options, sessionData, err := h.webAuthn.BeginRegistration(user)
	if err != nil {
		logger.Error(ctx, "Error beginning registration", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// セッション作成
	session, err := h.sr.Create(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to create session", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Username = req.Username
	session.RegistrationData = sessionData

	// Store に保存
	err = h.sr.Save(ctx, session)
	if err != nil {
		logger.Error(ctx, "Failed to store challenge", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// クッキー生成
	h.setSessionCookie(&w, session)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(options); err != nil {
		logger.Error(ctx, "Failed to write response", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *auth) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "finish registration ----------------------")

	var req request.FinishUserRegister
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info(ctx, "can't decode user data", logger.WithError(err))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		logger.Info(ctx, "username is empty")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// ユーザー確認
	exists, err := h.ur.ExistsByUsername(ctx, req.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		logger.Info(ctx, fmt.Sprintf("Exists User name: %s", req.Username))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// セッション確認
	cookie, err := r.Cookie("session")
	if err != nil {
		logger.Info(ctx, "session cookie is not found")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	session, err := h.sr.Get(ctx, cookie.Value)
	if err != nil {
		logger.Error(ctx, "can't get session", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if session == nil || session.RegistrationData == nil {
		logger.Info(ctx, "session is nil")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	var user model.User
	user.ID = cookie.Value
	user.Name = req.Username
	user.DisplayName = req.Username
	h.ur.Create(ctx, &user)
	return
}
