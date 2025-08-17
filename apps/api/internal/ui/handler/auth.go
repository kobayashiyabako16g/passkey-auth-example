package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/ui/handler/request"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/usecase"
	dtos "github.com/kobayashiyabako16g/passkey-auth-example/internal/usecase/dto/auth"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type Auth interface {
	BeginRegistration(w http.ResponseWriter, r *http.Request)
	FinishRegistration(w http.ResponseWriter, r *http.Request)
	BeginLogin(w http.ResponseWriter, r *http.Request)
	FinishLogin(w http.ResponseWriter, r *http.Request)
}

type auth struct {
	usecase usecase.Auth
}

func NewAuth(usecase usecase.Auth) Auth {
	return &auth{usecase}
}

func (s *auth) setSessionCookie(w *http.ResponseWriter, session *model.Session) {
	http.SetCookie(*w, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Required for cross-origin requests
		Secure:   true,                  // Set to true in production with HTTPS
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

	result, err := h.usecase.BeginRegistration(ctx, dtos.BeginRegistrationRequest{
		Username: req.Username,
	})
	if err != nil {
		switch err {
		case dtos.ErrUserExists:
			http.Error(w, "Username already exists", http.StatusConflict)
		default:
			logger.Error(ctx, "Failed to begin registration", logger.WithError(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// クッキー生成
	h.setSessionCookie(&w, result.Session)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.Cred); err != nil {
		logger.Error(ctx, "Failed to write response", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *auth) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "finish registration ----------------------")

	// セッション確認
	cookie, err := r.Cookie("session")
	if err != nil {
		logger.Info(ctx, "Handler: session cookie is not found")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	if cookie.Value == "" {
		logger.Info(ctx, "Handler: session cookie is empty")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	err = h.usecase.FinishRegistration(ctx, dtos.FinishRegistrationRequest{
		Session: cookie.Value,
		Request: r,
	})
	if err != nil {
		switch err {
		case dtos.ErrUserExists:
			http.Error(w, "Bad Requset", http.StatusBadRequest)
		case dtos.ErrFinishRegistration:
			http.Error(w, "Bad Requset", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		// clean up sid cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "sid",
			Value: "",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sid",
		Value: "",
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *auth) BeginLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "begin Login ----------------------")

	var req request.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info(ctx, "can't decode user data", logger.WithError(err))
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	var login dtos.BeginLoginRequest
	login.Username = req.Username
	// セッション確認
	cookie, err := r.Cookie("session")
	if err == nil {
		logger.Info(ctx, "session cookie is not found")
		login.Session = cookie.Value
	}

	// usecase
	result, err := h.usecase.BeginLogin(ctx, login)
	if err != nil {
		switch err {
		case dtos.ErrUserNotFound:
			http.Error(w, "Bad Requset", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	h.setSessionCookie(&w, result.Session)

	// option返却
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.Cred); err != nil {
		logger.Error(ctx, "Failed to write response", logger.WithError(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *auth) FinishLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "Finish login ----------------------")

	// セッション確認
	cookie, err := r.Cookie("session")
	if err != nil {
		logger.Info(ctx, "session cookie is not found")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}
	if cookie.Value == "" {
		logger.Info(ctx, "session cookie is empty")
		http.Error(w, "Bad Requset", http.StatusBadRequest)
		return
	}

	// usecase
	err = h.usecase.FinishLogin(ctx, dtos.FinishLoginRequest{
		Session: cookie.Value,
		Request: r,
	})
	if err != nil {
		switch err {
		case dtos.ErrSessionNotFound:
			http.Error(w, "Bad Requset", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
